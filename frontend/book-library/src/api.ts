import axios from 'axios'

const API_URL = 'http://localhost:8080/api'

export function setToken(token: string) {
	localStorage.setItem('token', token)
}

export function getToken() {
	return localStorage.getItem('token')
}

export function removeToken() {
	localStorage.removeItem('token')
}

const api = axios.create({
	baseURL: API_URL,
	withCredentials: true,
})

api.interceptors.request.use(config => {
	const token = getToken()
	console.log('Sending token:', token)
	if (token) {
		config.headers = config.headers || {}
		config.headers['Authorization'] = `Bearer ${token}`
	}
	return config
})

export async function register(email: string, password: string) {
	const res = await api.post('/register', { email, password })
	return res.data
}

export async function login(email: string, password: string) {
	const res = await api.post('/login', { email, password })
	setToken(res.data.JWT)
	return res.data
}

export async function getBooks() {
	const res = await api.get('/books')
	return res.data
}

export async function addBook(formData: FormData) {
	// POST /api/books/create (multipart)
	const token = getToken()
	const res = await api.post('/books/create', formData, {
		headers: token ? { Authorization: `Bearer ${token}` } : undefined,
	})
	return res.data
}

export async function deleteBook(id: number) {
	const res = await api.delete(`/books/delete?id=${id}`)
	return res.data
}

export async function getCategoriesWithBooks() {
	const res = await api.get('/categories/with-books')
	return res.data
}

export async function createCategory(name: string) {
	return api.post('/categories/create', { name }, { withCredentials: true })
}

export async function renameCategory(id: number, name: string) {
	return api.put(
		`/categories/rename?id=${id}`,
		{ name },
		{ withCredentials: true }
	)
}

export async function assignBookToCategory(bookId: number, categoryId: number) {
	const res = await api.post('/books/assign', { bookId, categoryId })
	return res.data
}

export async function getUncategorizedBooks() {
	const res = await api.get('/books/uncategorized')
	return res.data
}

export async function deleteCategory(id: number) {
	return api.delete(`/categories/delete?id=${id}`, { withCredentials: true })
}
