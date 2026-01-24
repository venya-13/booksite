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
			setError('Не удалось загрузить книги')
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
			setError(err?.response?.data?.error || 'Не удалось добавить книгу')
		} finally {
			setAdding(false)
		}
	}

	const handleRemoveBook = async (id: number) => {
		try {
			await deleteBook(id)
			await fetchBooks()
		} catch (e: any) {
			setError('Не удалось удалить книгу')
		}
	}

	return (
		<Box sx={{ px: { xs: 2, sm: 3, md: 4 }, py: 4, maxWidth: '1400px', mx: 'auto' }}>
			<Typography 
				variant='h4' 
				gutterBottom 
				sx={{ 
					fontWeight: 800,
					mb: 4,
					background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
					WebkitBackgroundClip: 'text',
					WebkitTextFillColor: 'transparent',
					backgroundClip: 'text',
					letterSpacing: '-0.02em',
				}}
			>
				Панель Администратора
			</Typography>
			<Typography 
				variant='h6' 
				gutterBottom 
				sx={{ 
					fontWeight: 600,
					mb: 3,
					color: 'text.primary',
				}}
			>
				Добавить новую книгу
			</Typography>
			<Box
				component='form'
				onSubmit={handleAddBook}
				sx={{
					maxWidth: 600,
					mb: 6,
					p: 4,
					borderRadius: 4,
					backgroundColor: 'background.paper',
					boxShadow: '0 8px 30px rgba(2,6,23,0.1)',
					transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
					'&:hover': {
						boxShadow: '0 12px 40px rgba(2,6,23,0.12)',
					},
				}}
			>
				<Stack spacing={3}>
					<TextField
						label='Название'
						name='title'
						value={form.title}
						onChange={handleChange}
						required
						fullWidth
						sx={{
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
					/>
					<TextField
						label='Автор'
						name='author'
						value={form.author}
						onChange={handleChange}
						required
						fullWidth
						sx={{
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
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
						fullWidth
						label='Категория'
						sx={{
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
					>
						<option value=''>Без категории</option>

						{useCategoriesStore.getState().categories.map(cat => (
							<option key={cat.id} value={cat.id}>
								{cat.is_system ? `${cat.name}` : cat.name}
							</option>
						))}
					</TextField>
					<TextField
						label='URL книги'
						name='bookUrl'
						value={form.bookUrl || ''}
						onChange={handleChange}
						required
						fullWidth
						sx={{
							'& .MuiOutlinedInput-root': {
								borderRadius: 2,
							},
						}}
					/>
					<Box>
						<Button
							component='label'
							variant='outlined'
							fullWidth
							sx={{
								py: 1.5,
								borderRadius: 2,
								borderWidth: 2,
								'&:hover': {
									borderWidth: 2,
								},
							}}
						>
							Выбрать обложку
							<input
								type='file'
								accept='.jpg,.jpeg,.png'
								onChange={handleFileChange}
								hidden
							/>
						</Button>

						{form.cover && typeof form.cover !== 'string' && (
							<Box
								sx={{
									mt: 2,
									p: 2,
									borderRadius: 2,
									backgroundColor: 'rgba(25, 118, 210, 0.04)',
									display: 'inline-block',
								}}
							>
								<img
									src={URL.createObjectURL(form.cover)}
									alt='Превью обложки книги'
									style={{
										width: 200,
										borderRadius: 12,
										border: '2px solid rgba(25, 118, 210, 0.2)',
										boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
									}}
								/>
							</Box>
						)}
					</Box>

					{error && (
						<Box
							sx={{
								p: 2,
								borderRadius: 2,
								backgroundColor: 'error.light',
								color: 'error.main',
							}}
						>
							<Typography variant='body2'>{error}</Typography>
						</Box>
					)}

					<Button 
						type='submit' 
						variant='contained' 
						disabled={adding}
						fullWidth
						size='large'
						sx={{
							py: 1.5,
							borderRadius: 2,
							background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
							'&:hover': {
								background: 'linear-gradient(135deg, #1565c0 0%, #6a1b9a 100%)',
							},
							'&:disabled': {
								background: 'rgba(0,0,0,0.12)',
							},
						}}
					>
						{adding ? 'Добавление...' : 'Добавить книгу'}
					</Button>
				</Stack>
			</Box>
			<Typography 
				variant='h6' 
				gutterBottom 
				sx={{ 
					fontWeight: 600,
					mb: 3,
					color: 'text.primary',
				}}
			>
				Текущие книги
			</Typography>
			{loading ? (
				<Box display='flex' justifyContent='center' mt={6} mb={6}>
					<CircularProgress size={60} thickness={4} />
				</Box>
			) : error ? (
				<Box
					sx={{
						p: 3,
						borderRadius: 2,
						backgroundColor: 'error.light',
						color: 'error.main',
					}}
				>
					<Typography color='error'>{error}</Typography>
				</Box>
			) : (
				<Box 
					display='flex' 
					flexWrap='wrap' 
					gap={3}
					justifyContent={{ xs: 'center', md: 'flex-start' }}
				>
					{books.map(book => (
						<Card
							key={book.id}
							sx={{
								width: { xs: '100%', sm: 280, md: 300 },
								maxWidth: 300,
								overflow: 'hidden',
							}}
						>
							{book.cover_path && (
								<Box
									component='img'
									src={`http://localhost:8080${book.cover_path}`}
									alt={book.title}
									sx={{
										width: '100%',
										height: 220,
										objectFit: 'cover',
										transition: 'transform 0.5s cubic-bezier(0.4, 0, 0.2, 1)',
										'&:hover': {
											transform: 'scale(1.05)',
										},
									}}
								/>
							)}

							<CardContent sx={{ pb: 2, pt: 2.5 }}>
								<Typography 
									variant='h6' 
									sx={{ 
										fontWeight: 700,
										mb: 1,
										lineHeight: 1.3,
										minHeight: '3em',
										display: '-webkit-box',
										WebkitLineClamp: 2,
										WebkitBoxOrient: 'vertical',
										overflow: 'hidden',
									}}
								>
									{book.title}
								</Typography>
								<Typography 
									color='text.secondary'
									variant='body2'
									sx={{ 
										fontWeight: 500,
										opacity: 0.8,
									}}
								>
									от {book.author}
								</Typography>
							</CardContent>
							<CardActions sx={{ px: 2, pb: 2.5 }}>
								<Button 
									color='error' 
									onClick={() => handleRemoveBook(book.id)}
									fullWidth
									variant='outlined'
									sx={{
										borderWidth: 2,
										borderRadius: 2,
										py: 1,
										'&:hover': {
											borderWidth: 2,
											backgroundColor: 'error.light',
										},
									}}
								>
									Удалить
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
