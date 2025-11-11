package httpserver

import (
	"encoding/json"
	"fmt"
	"google-auth-demo/backend/internal/repo"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		cats := r.FormValue("categories")
		fmt.Println("CATEGORIES RAW:", cats)

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

func (s *Server) handleGetCategoriesWithBooks(w http.ResponseWriter, r *http.Request) {
	// call service that returns categories with books
	cats, err := s.svc.GetCategoriesWithBooks(r.Context())
	if err != nil {
		http.Error(w, "failed to get categories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

func (s *Server) handleGetUncategorizedBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.svc.GetUncategorizedBooks(r.Context())
	if err != nil {
		http.Error(w, "failed to get uncategorized books: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (s *Server) handleAssignBookToCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		BookID     int `json:"bookId"`
		CategoryID int `json:"categoryId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if body.BookID == 0 || body.CategoryID == 0 {
		http.Error(w, "bookId and categoryId required", http.StatusBadRequest)
		return
	}
	if err := s.svc.AssignBookToCategory(r.Context(), body.BookID, body.CategoryID); err != nil {
		http.Error(w, "failed to assign: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
