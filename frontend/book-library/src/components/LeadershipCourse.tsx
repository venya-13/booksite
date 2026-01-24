import React, { useEffect, useMemo, useState } from 'react'
import {
	Box,
	Typography,
	Card,
	CardContent,
	CardActions,
	CardMedia,
	Button,
	CircularProgress,
	Divider,
} from '@mui/material'
import { getCategoriesWithBooks } from '../api'

const placeholderCover = 'https://via.placeholder.com/300x180?text=Book+Cover'
const LEADERSHIP_SLUG = 'leadership-course'

const LeadershipCourse: React.FC = () => {
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<string | null>(null)
	const [category, setCategory] = useState<any | null>(null)

	useEffect(() => {
		setLoading(true)
		setError(null)
		getCategoriesWithBooks()
			.then(data => (Array.isArray(data) ? data : []))
			.then(cats => cats.find((c: any) => c?.slug === LEADERSHIP_SLUG) || null)
			.then(cat => setCategory(cat))
			.catch(() => setError('Не удалось загрузить курс'))
			.finally(() => setLoading(false))
	}, [])

	const books = useMemo(() => (category?.books ? category.books : []), [category])

	return (
		<Box p={2}>
			<Typography variant='h4' gutterBottom sx={{ fontWeight: 800 }}>
				Курс Лидерство
			</Typography>

			{loading ? (
				<Box display='flex' justifyContent='center' mt={4}>
					<CircularProgress />
				</Box>
			) : error ? (
				<Typography color='error'>{error}</Typography>
			) : !category ? (
				<Typography>
					Категория курса не найдена (ожидается slug: {LEADERSHIP_SLUG})
				</Typography>
			) : (
				<Box
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
							{category.name}
						</Typography>
					</Divider>
					{books.length === 0 ? (
						<Typography sx={{ px: 2, pb: 2 }}>
							Пока нет книг в этом курсе.
						</Typography>
					) : (
						<Box
							display='flex'
							flexWrap='wrap'
							gap={2}
							justifyContent='center'
						>
							{books.map((book: any) => (
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
										alt={book.title + ' обложка'}
									/>
									<CardContent>
										<Typography variant='h6' sx={{ fontWeight: 700 }}>
											{book.title}
										</Typography>
										<Typography color='text.secondary'>от {book.author}</Typography>
									</CardContent>
									<CardActions>
										<Button
											size='small'
											onClick={() => window.open(book.file_url, '_blank')}
										>
											Читать
										</Button>
									</CardActions>
								</Card>
							))}
						</Box>
					)}
				</Box>
			)}
		</Box>
	)
}

export default LeadershipCourse

