import React, { useState, useEffect } from 'react'
import {
	Box,
	Typography,
	TextField,
	Card,
	CardContent,
	CardActions,
	Button,
	Dialog,
	DialogTitle,
	DialogContent,
	CardMedia,
	CircularProgress,
	Divider,
	InputAdornment,
	IconButton,
} from '@mui/material'
import { Search, Clear } from '@mui/icons-material'
import { getCategoriesWithBooks, getUncategorizedBooks } from '../api'

const placeholderCover = 'https://via.placeholder.com/300x180?text=Book+Cover'

const Books: React.FC = () => {
	const [categories, setCategories] = useState<any[]>([])
	const [search, setSearch] = useState('')
	const [selectedBook, setSelectedBook] = useState<any | null>(null)
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<string | null>(null)
	const [uncategorizedBooks, setUncategorizedBooks] = useState<any[]>([])

	useEffect(() => {
		setLoading(true)
		setError(null)
		Promise.all([
			getCategoriesWithBooks().then(data => (Array.isArray(data) ? data : [])),
			getUncategorizedBooks().then(data => (Array.isArray(data) ? data : [])),
		])
			.then(([cats, uncats]) => {
				setCategories(cats)
				setUncategorizedBooks(uncats)
			})
			.catch(() => setError('Failed to load books'))
			.finally(() => setLoading(false))
	}, [])

	const filteredCategories = categories
		.map((cat: any) => ({
			...cat,
			books: (cat.books || []).filter(
				(book: any) =>
					book.title.toLowerCase().includes(search.toLowerCase()) ||
					book.author.toLowerCase().includes(search.toLowerCase())
			),
		}))
		.filter(cat => cat.books.length > 0)

	const filteredUncategorizedBooks = (uncategorizedBooks || []).filter(
		book =>
			book.title.toLowerCase().includes(search.toLowerCase()) ||
			book.author.toLowerCase().includes(search.toLowerCase())
	)

	const handleClearSearch = () => setSearch('')

	return (
		<Box p={2}>
			<Typography variant='h4' gutterBottom sx={{ fontWeight: 800 }}>
				Browse and Read Books
			</Typography>

			<TextField
				label='Search books'
				variant='outlined'
				fullWidth
				margin='normal'
				value={search}
				onChange={e => setSearch(e.target.value)}
				InputProps={{
					startAdornment: (
						<InputAdornment position='start'>
							<Search />
						</InputAdornment>
					),
					endAdornment: search && (
						<InputAdornment position='end'>
							<IconButton onClick={handleClearSearch} edge='end'>
								<Clear />
							</IconButton>
						</InputAdornment>
					),
				}}
				sx={{
					'& .MuiOutlinedInput-root': {
						backgroundColor: 'background.paper',
						borderRadius: 3,
						boxShadow: '0 6px 18px rgba(2, 6, 23, 0.06)',
					},
				}}
			/>

			{loading ? (
				<Box display='flex' justifyContent='center' mt={4}>
					<CircularProgress />
				</Box>
			) : error ? (
				<Typography color='error'>{error}</Typography>
			) : filteredCategories.length === 0 &&
			  filteredUncategorizedBooks.length === 0 ? (
				<Typography>No books found.</Typography>
			) : (
				<>
					{filteredCategories.map((cat: any) => (
						<Box
							key={cat.id}
							mb={5}
							sx={{
								p: 1,
								borderRadius: 3,
								background:
									'linear-gradient(180deg, rgba(25,118,210,0.06), rgba(124,77,255,0.04))',
							}}
						>
							<Divider sx={{ mb: 2 }}>
								<Typography variant='h5' color='primary'>
									{cat.name}
								</Typography>
							</Divider>
							<Box
								display='flex'
								flexWrap='wrap'
								gap={2}
								justifyContent='center'
							>
								{cat.books.map((book: any) => (
									<Card
										key={book.id}
										sx={{
											width: 300,
											flex: '0 1 300px',
											'&:hover': {
												transform: 'translateY(-4px)',
												boxShadow: '0 16px 30px rgba(2,6,23,0.10)',
											},
										}}
									>
										<CardMedia
											component='img'
											height='180'
											image={
												book.cover_path
													? book.cover_path.startsWith('/')
														? `http://localhost:8080${book.cover_path}`
														: book.cover_path
													: placeholderCover
											}
											alt={book.title + ' cover'}
										/>
										<CardContent>
											<Typography variant='h6' sx={{ fontWeight: 700 }}>
												{book.title}
											</Typography>
											<Typography color='text.secondary'>
												by {book.author}
											</Typography>
										</CardContent>
										<CardActions>
											<Button
												size='small'
												onClick={() => setSelectedBook(book)}
											>
												Read
											</Button>
										</CardActions>
									</Card>
								))}
							</Box>
						</Box>
					))}

					{filteredUncategorizedBooks.length > 0 && (
						<Box mb={5}>
							<Divider sx={{ mb: 2 }}>
								<Typography variant='h5' color='secondary'>
									Uncategorized
								</Typography>
							</Divider>
							<Box
								display='flex'
								flexWrap='wrap'
								gap={2}
								justifyContent='center'
							>
								{filteredUncategorizedBooks.map((book: any) => (
									<Card
										key={book.id}
										sx={{
											width: 300,
											flex: '0 1 300px',
											'&:hover': {
												transform: 'translateY(-4px)',
												boxShadow: '0 16px 30px rgba(2,6,23,0.10)',
											},
										}}
									>
										<CardMedia
											component='img'
											height='180'
											image={
												book.cover_path
													? book.cover_path.startsWith('/')
														? `http://localhost:8080${book.cover_path}`
														: book.cover_path
													: placeholderCover
											}
											alt={book.title + ' cover'}
										/>
										<CardContent>
											<Typography variant='h6' sx={{ fontWeight: 700 }}>
												{book.title}
											</Typography>
											<Typography color='text.secondary'>
												by {book.author}
											</Typography>
										</CardContent>
										<CardActions>
											<Button
												size='small'
												onClick={() => setSelectedBook(book)}
											>
												Read
											</Button>
										</CardActions>
									</Card>
								))}
							</Box>
						</Box>
					)}
				</>
			)}

			<Dialog
				open={!!selectedBook}
				onClose={() => setSelectedBook(null)}
				maxWidth='sm'
				fullWidth
			>
				<DialogTitle>{selectedBook?.title}</DialogTitle>
				<DialogContent>
					<Typography variant='subtitle1' gutterBottom>
						by {selectedBook?.author}
					</Typography>
					<Button
						variant='contained'
						href={selectedBook?.file_url}
						target='_blank'
					>
						Open Book
					</Button>
				</DialogContent>
			</Dialog>
		</Box>
	)
}

export default Books
