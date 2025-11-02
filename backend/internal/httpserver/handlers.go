package httpserver

import (
	"encoding/json"
	"fmt"
	"google-auth-demo/backend/internal/repo"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Google Auth Demo!")
	fmt.Fprintln(w, "Click here to <a href='/login'>Login with Google</a>")
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := s.svc.GetAuthURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Warn("Missing code in callback", slog.String("url", r.URL.String()))
		http.Error(w, "No code found in callback", http.StatusBadRequest)
		return
	}

	slog.Info("OAuth callback received", slog.String("code", code[:6]+"...")) // log only first 6 chars for privacy

	resp, err := s.svc.HandleCallback(code)
	if err != nil {
		slog.Error("HandleCallback failed", slog.String("error", err.Error()))
		http.Error(w, "Callback error", http.StatusInternalServerError)
		return
	}

	id, _ := resp.User["id"].(string)
	email, _ := resp.User["email"].(string)
	slog.Info("User authenticated successfully",
		slog.String("google_id", id),
		slog.String("email", email),
	)

	redirectURL := s.svc.GetFrontendURL(resp.JWT)
	slog.Info("Redirecting user to frontend", slog.String("url", redirectURL))

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (s *Server) handleGoogleProfile(w http.ResponseWriter, r *http.Request) {
	googleID, _ := r.Context().Value("user_id").(string)
	if googleID == "" {
		http.Error(w, "missing user_id in context", http.StatusUnauthorized)
		return
	}

	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := s.svc.FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "Google API error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// return JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	googleID := r.URL.Query().Get("google_id")
	if googleID == "" {
		http.Error(w, "google_id required", http.StatusBadRequest)
		return
	}

	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "failed to refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := s.svc.FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "failed to fetch profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo["id"] = googleID
	userInfo["access_token"] = accessToken
	if err := s.svc.SaveUser(userInfo); err != nil {
		http.Error(w, "failed to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token": accessToken,
		"user":         userInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	googleID := r.URL.Query().Get("google_id")
	if googleID == "" {
		http.Error(w, "google_id required", http.StatusBadRequest)
		return
	}

	// gет access token
	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "failed to ensure token", http.StatusInternalServerError)
		return
	}

	// get profile
	profile, err := s.svc.FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "failed to fetch profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}

func (s *Server) handleProtected(w http.ResponseWriter, r *http.Request) {
	// values come from middleware
	userID, _ := r.Context().Value("user_id").(string)
	email, _ := r.Context().Value("email").(string)
	isAdmin, _ := r.Context().Value("is_admin").(bool)

	resp := map[string]interface{}{
		"message":  "This is a protected route",
		"user_id":  userID,
		"email":    email,
		"is_admin": isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleJWTRefresh(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		slog.Warn("Refresh token missing",
			slog.String("remote_ip", r.RemoteAddr),
			slog.String("path", r.URL.Path),
		)
		http.Error(w, "refresh_token missing", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	resp, err := s.svc.RefreshJWT(refreshToken)
	if err != nil {
		slog.Error("Failed to refresh JWT",
			slog.String("remote_ip", r.RemoteAddr),
			slog.String("error", err.Error()),
			slog.String("path", r.URL.Path),
			slog.Duration("elapsed", time.Since(start)),
		)
		http.Error(w, "invalid refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.svc.JWTTTL.Seconds()),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.svc.RefreshTTL.Seconds()),
	})

	slog.Info("JWT successfully refreshed",
		slog.String("remote_ip", r.RemoteAddr),
		slog.String("path", r.URL.Path),
		slog.Duration("elapsed", time.Since(start)),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "JWT refreshed successfully"}`))
}

func (s *Server) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := s.svc.GetAllCategories(r.Context())
	if err != nil {
		http.Error(w, "failed to get categories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (s *Server) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := s.svc.CreateCategory(r.Context(), req.Name); err != nil {
		http.Error(w, "failed to create category: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := s.svc.DeleteCategory(r.Context(), id); err != nil {
		http.Error(w, "failed to delete category: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleRenameCategory(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing category id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := s.svc.RenameCategory(ctx, id, body.Name); err != nil {
		http.Error(w, "failed to rename category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// limit and parse
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	description := r.FormValue("description")
	fileURL := r.FormValue("file_url")
	cats := r.FormValue("categories")

	var categoryIDs []int
	for _, sID := range strings.Split(cats, ",") {
		if sID == "" {
			continue
		}
		id, _ := strconv.Atoi(sID)
		categoryIDs = append(categoryIDs, id)
	}

	coverPath := ""
	file, header, err := r.FormFile("cover")
	if err == nil && file != nil {
		defer file.Close()
		os.MkdirAll("./uploads/covers", 0755)
		// unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
		dstPath := "./uploads/covers/" + filename
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "failed to write file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		coverPath = "/covers/" + filename // path for frontend
	}

	book := repo.Book{
		Title:       title,
		Author:      author,
		Description: description,
		FileURL:     fileURL,
		CoverPath:   coverPath,
	}

	id, err := s.svc.CreateBook(r.Context(), book, categoryIDs)
	if err != nil {
		http.Error(w, "create failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (s *Server) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.svc.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, "failed to get books: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (s *Server) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	book, err := s.svc.GetBookByID(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get book: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	title := r.FormValue("title")
	author := r.FormValue("author")
	description := r.FormValue("description")
	fileURL := r.FormValue("file_url")

	cats := r.FormValue("categories")
	var categoryIDs []int
	for _, sID := range strings.Split(cats, ",") {
		if sID == "" {
			continue
		}
		id, _ := strconv.Atoi(sID)
		categoryIDs = append(categoryIDs, id)
	}

	coverPath := r.FormValue("existing_cover")

	file, header, err := r.FormFile("cover")
	if err == nil && file != nil {
		defer file.Close()
		os.MkdirAll("./uploads/covers", 0755)
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
		dstPath := "./uploads/covers/" + filename
		dst, _ := os.Create(dstPath)
		io.Copy(dst, file)
		coverPath = "/covers/" + filename
	}

	book := repo.Book{
		ID:          id,
		Title:       title,
		Author:      author,
		Description: description,
		FileURL:     fileURL,
		CoverPath:   coverPath,
	}

	if err := s.svc.UpdateBook(r.Context(), book, categoryIDs); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("updated"))
}

func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := s.svc.DeleteBook(r.Context(), id); err != nil {
		http.Error(w, "delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("deleted"))
}

func (s *Server) handleGetBooksByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category"))
	books, err := s.svc.GetBooksByCategory(r.Context(), categoryID)
	if err != nil {
		http.Error(w, "failed to get books: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (s *Server) UploadBookCover(w http.ResponseWriter, r *http.Request) {
	// limit form size to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// read file
	file, header, err := r.FormFile("cover")
	if err != nil {
		http.Error(w, "cover file is required: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// get book id from query param
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing book id query param", http.StatusBadRequest)
		return
	}
	bookID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book id: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ensure uploads directory exists
	if err := os.MkdirAll("./uploads/covers", 0o755); err != nil {
		http.Error(w, "failed to create uploads dir: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	dstPath := "./uploads/covers/" + filename

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// build public path for frontend (served from /covers/)
	publicPath := "/covers/" + filename

	// load existing book (so we don't wipe other fields or categories)
	book, err := s.svc.GetBookByID(r.Context(), bookID)
	if err != nil {
		http.Error(w, "failed to get book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// update cover path on book struct
	book.CoverPath = publicPath

	// call service update -- keep categories unchanged (pass nil or empty slice depending on your repo impl)
	if err := s.svc.UpdateBook(r.Context(), *book, nil); err != nil {
		http.Error(w, "failed to update book record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"cover_url":"%s"}`, publicPath)))
}
