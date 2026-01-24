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
			.catch(() => setError('Не удалось загрузить книги'))
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
				Просматривайте и читайте книги
			</Typography>

			<TextField
				label='Поиск книг'
				variant='outlined'
				fullWidth
				margin='normal'
				value={search}
				onChange={e => setSearch(e.target.value)}
				InputProps={{
					startAdornment: (
						<InputAdornment position='start'>
							<Search sx={{ color: 'primary.main', opacity: 0.7 }} />
						</InputAdornment>
					),
					endAdornment: search && (
						<InputAdornment position='end'>
							<IconButton 
								onClick={handleClearSearch} 
								edge='end'
								sx={{
									'&:hover': {
										backgroundColor: 'rgba(25, 118, 210, 0.08)',
									},
								}}
							>
								<Clear />
							</IconButton>
						</InputAdornment>
					),
				}}
				sx={{
					mb: 4,
					'& .MuiOutlinedInput-root': {
						backgroundColor: 'background.paper',
						borderRadius: 4,
						boxShadow: '0 4px 16px rgba(2, 6, 23, 0.08)',
						fontSize: '1.05rem',
						'& fieldset': {
							borderWidth: '2px',
						},
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
				<Typography>Книг не найдено.</Typography>
			) : (
				<>
					{filteredCategories.map((cat: any) => (
						<Box
							key={cat.id}
							mb={6}
							sx={{
								p: 3,
								borderRadius: 4,
								background:
									'linear-gradient(180deg, rgba(25,118,210,0.08), rgba(124,77,255,0.06))',
								boxShadow: '0 4px 20px rgba(25, 118, 210, 0.08)',
								transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
								'&:hover': {
									boxShadow: '0 8px 30px rgba(25, 118, 210, 0.12)',
								},
							}}
						>
							<Divider sx={{ mb: 4 }}>
								<Typography 
									variant='h5' 
									color='primary'
									sx={{ 
										fontWeight: 700,
										px: 3,
										py: 1,
										background: 'rgba(255, 255, 255, 0.9)',
										borderRadius: 2,
									}}
								>
									{cat.name}
								</Typography>
							</Divider>
							<Box
								display='flex'
								flexWrap='wrap'
								gap={3}
								justifyContent={{ xs: 'center', md: 'flex-start' }}
							>
								{cat.books.map((book: any) => (
									<Card
										key={book.id}
										sx={{
											width: { xs: '100%', sm: 280, md: 300 },
											maxWidth: 300,
											overflow: 'hidden',
											cursor: 'pointer',
										}}
									>
										<CardMedia
											component='img'
											height='220'
											image={
												book.cover_path
													? book.cover_path.startsWith('/')
														? `http://localhost:8080${book.cover_path}`
														: book.cover_path
													: placeholderCover
											}
											alt={book.title + ' обложка'}
											sx={{
												transition: 'transform 0.5s cubic-bezier(0.4, 0, 0.2, 1)',
												'&:hover': {
													transform: 'scale(1.05)',
												},
											}}
										/>
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
												size='medium'
												variant='contained'
												fullWidth
												onClick={() => window.open(book.file_url, '_blank')}
												sx={{
													py: 1.5,
													borderRadius: 2,
													background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
													'&:hover': {
														background: 'linear-gradient(135deg, #1565c0 0%, #6a1b9a 100%)',
													},
												}}
											>
												Читать
											</Button>
										</CardActions>
									</Card>
								))}
							</Box>
						</Box>
					))}

					{filteredUncategorizedBooks.length > 0 && (
						<Box 
							mb={6}
							sx={{
								p: 3,
								borderRadius: 4,
								background:
									'linear-gradient(180deg, rgba(124,77,255,0.08), rgba(25,118,210,0.06))',
								boxShadow: '0 4px 20px rgba(124, 77, 255, 0.08)',
							}}
						>
							<Divider sx={{ mb: 4 }}>
								<Typography 
									variant='h5' 
									color='secondary'
									sx={{ 
										fontWeight: 700,
										px: 3,
										py: 1,
										background: 'rgba(255, 255, 255, 0.9)',
										borderRadius: 2,
									}}
								>
									Без категории
								</Typography>
							</Divider>
							<Box
								display='flex'
								flexWrap='wrap'
								gap={3}
								justifyContent={{ xs: 'center', md: 'flex-start' }}
							>
								{filteredUncategorizedBooks.map((book: any) => (
									<Card
										key={book.id}
										sx={{
											width: { xs: '100%', sm: 280, md: 300 },
											maxWidth: 300,
											overflow: 'hidden',
											cursor: 'pointer',
										}}
									>
										<CardMedia
											component='img'
											height='220'
											image={
												book.cover_path
													? book.cover_path.startsWith('/')
														? `http://localhost:8080${book.cover_path}`
														: book.cover_path
													: placeholderCover
											}
											alt={book.title + ' обложка'}
											sx={{
												transition: 'transform 0.5s cubic-bezier(0.4, 0, 0.2, 1)',
												'&:hover': {
													transform: 'scale(1.05)',
												},
											}}
										/>
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
												size='medium'
												variant='contained'
												fullWidth
												onClick={() => window.open(book.file_url, '_blank')}
												sx={{
													py: 1.5,
													borderRadius: 2,
													background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
													'&:hover': {
														background: 'linear-gradient(135deg, #1565c0 0%, #6a1b9a 100%)',
													},
												}}
											>
												Читать
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
						от {selectedBook?.author}
					</Typography>
					<Button
						variant='contained'
						href={selectedBook?.file_url}
						target='_blank'
					>
						Открыть книгу
					</Button>
				</DialogContent>
			</Dialog>
		</Box>
	)
}

export default Books
