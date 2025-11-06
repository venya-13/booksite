import React, { useState, useEffect } from 'react'
import { useCategoriesStore } from '../store/categoriesStore'
import {
	Box,
	Typography,
	TextField,
	Button,
	Stack,
	Card,
	CardContent,
	CardActions,
	CircularProgress,
} from '@mui/material'
import {
	getBooks,
	addBook,
	deleteBook,
	getCategoriesWithBooks,
	assignBookToCategory,
	getUncategorizedBooks,
	createCategory,
	renameCategory,
	deleteCategory,
	getUncategorizedBooks as getUncatBooks,
	getCategoriesWithBooks as getCatBooks,
} from '../api'
import axios from 'axios'

const AdminPanel: React.FC = () => {
	const [books, setBooks] = useState<any[]>([])
	const [form, setForm] = useState({
		title: '',
		author: '',
		bookUrl: '',
		cover: null as File | null,
	})

	const refreshCategories = useCategoriesStore(state => state.refresh)

	const [selectedCategory, setSelectedCategory] = useState<number | ''>('')

	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<string | null>(null)
	const [adding, setAdding] = useState(false)

	const fetchBooks = async () => {
		setLoading(true)
		setError(null)
		try {
			const data = await getBooks()
			setBooks(Array.isArray(data) ? data : [])
		} catch (e: any) {
			setError('Failed to load books')
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => {
		fetchBooks()
		refreshCategories()
	}, [])

	const handleChange = (
		e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => {
		setForm({ ...form, [e.target.name]: e.target.value })
	}

	const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		if (e.target.files && e.target.files[0]) {
			const file = e.target.files[0]
			setForm(prev => ({ ...prev, cover: file }))
		}
	}

	const handleAddBook = async (e: React.FormEvent) => {
		e.preventDefault()

		if (!form.title || !form.author || !form.bookUrl) return

		const fd = new FormData()
		fd.append('title', form.title)
		fd.append('author', form.author)
		fd.append('file_url', form.bookUrl)
		if (form.cover) fd.append('cover', form.cover)

		if (selectedCategory !== '') {
			fd.append('categories', String(selectedCategory))
		}

		setAdding(true)
		setError(null)

		try {
			const token = localStorage.getItem('token')

			const createRes = await axios.post(
				'http://localhost:8080/api/books/create',
				fd,
				{
					withCredentials: true,
					headers: {
						Authorization: token ? `Bearer ${token}` : '',
					},
				}
			)

			const newBookId = createRes.data?.id
			console.log('Book created:', newBookId)

			if (selectedCategory !== '' && newBookId) {
				await axios.post(
					'http://localhost:8080/api/books/assign',
					{
						bookId: newBookId,
						categoryId: Number(selectedCategory),
					},
					{
						withCredentials: true,
						headers: {
							Authorization: token ? `Bearer ${token}` : '',
							'Content-Type': 'application/json',
						},
					}
				)
				console.log(
					`Book ${newBookId} assigned to category ${selectedCategory}`
				)
			}

			setForm({ title: '', author: '', bookUrl: '', cover: null })
			setSelectedCategory('')
			await fetchBooks()
			await refreshCategories()
		} catch (err: any) {
			console.error(
				'Add book error:',
				err?.response?.data ?? err.message ?? err
			)
			setError(err?.response?.data?.error || 'Failed to add book')
		} finally {
			setAdding(false)
		}
	}

	const handleRemoveBook = async (id: number) => {
		try {
			await deleteBook(id)
			await fetchBooks()
		} catch (e: any) {
			setError('Failed to remove book')
		}
	}

	return (
		<Box p={2}>
			<Typography variant='h4' gutterBottom sx={{ fontWeight: 800 }}>
				Admin Panel
			</Typography>
			<Typography variant='h6' gutterBottom>
				Add a New Book
			</Typography>
			<Box
				component='form'
				onSubmit={handleAddBook}
				sx={{
					maxWidth: 520,
					mb: 4,
					p: 3,
					borderRadius: 3,
					backgroundColor: 'background.paper',
					boxShadow: '0 10px 25px rgba(2,6,23,0.06)',
				}}
			>
				<Stack spacing={2}>
					<TextField
						label='Title'
						name='title'
						value={form.title}
						onChange={handleChange}
						required
					/>
					<TextField
						label='Author'
						name='author'
						value={form.author}
						onChange={handleChange}
						required
					/>

					<TextField
						select
						value={selectedCategory}
						onChange={e =>
							setSelectedCategory(
								e.target.value === '' ? '' : Number(e.target.value)
							)
						}
						SelectProps={{ native: true }}
					>
						<option value=''>No Category</option>

						{useCategoriesStore.getState().categories.map(cat => (
							<option key={cat.id} value={cat.id}>
								{cat.name}
							</option>
						))}
					</TextField>
					<TextField
						label='Book URL'
						name='bookUrl'
						value={form.bookUrl || ''}
						onChange={handleChange}
						required
					/>
					<>
						<input
							type='file'
							accept='.jpg,.jpeg,.png'
							onChange={handleFileChange}
							style={{ marginTop: 8 }}
						/>

						{form.cover && typeof form.cover !== 'string' && (
							<img
								src={URL.createObjectURL(form.cover)}
								alt='Book cover preview'
								style={{
									width: 160,
									marginTop: 8,
									borderRadius: 8,
									border: '1px solid #e5e7eb',
								}}
							/>
						)}
					</>

					<Button type='submit' variant='contained' disabled={adding}>
						{adding ? 'Adding...' : 'Add Book'}
					</Button>
				</Stack>
			</Box>
			<Typography variant='h6' gutterBottom>
				Current Books
			</Typography>
			{loading ? (
				<Box display='flex' justifyContent='center' mt={4}>
					<CircularProgress />
				</Box>
			) : error ? (
				<Typography color='error'>{error}</Typography>
			) : (
				<Box display='flex' flexWrap='wrap' gap={2}>
					{books.map(book => (
						<Card
							key={book.id}
							sx={{
								width: 300,
								'&:hover': {
									transform: 'translateY(-4px)',
									boxShadow: '0 16px 30px rgba(2,6,23,0.10)',
								},
							}}
						>
							{book.cover_path && (
								<img
									src={`http://localhost:8080${book.cover_path}`}
									alt={book.title}
									style={{ width: '100%', height: 180, objectFit: 'cover' }}
								/>
							)}

							<CardContent>
								<Typography variant='h6' sx={{ fontWeight: 700 }}>
									{book.title}
								</Typography>
								<Typography color='text.secondary'>by {book.author}</Typography>
							</CardContent>
							<CardActions>
								<Button color='error' onClick={() => handleRemoveBook(book.id)}>
									Remove
								</Button>
							</CardActions>
						</Card>
					))}
				</Box>
			)}
		</Box>
	)
}

export default AdminPanel
