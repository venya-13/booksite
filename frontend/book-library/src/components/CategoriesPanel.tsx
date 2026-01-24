import React, { useEffect, useState } from 'react'
import { useCategoriesStore } from '../store/categoriesStore'
import {
	getCategoriesWithBooks,
	renameCategory,
	createCategory,
	assignBookToCategory,
	getUncategorizedBooks,
	deleteCategory,
} from '../api'
import {
	Box,
	Typography,
	TextField,
	Button,
	Stack,
	Card,
	Dialog,
	DialogTitle,
	DialogContent,
	DialogActions,
	MenuItem,
	Select,
	InputLabel,
	FormControl,
} from '@mui/material'

const CategoriesPanel: React.FC = () => {
	const { categories, uncategorizedBooks, refresh } = useCategoriesStore()
	const [editId, setEditId] = useState<number | null>(null)
	const [editName, setEditName] = useState('')
	const [newCategory, setNewCategory] = useState('')
	const [assignDialog, setAssignDialog] = useState<{
		open: boolean
		bookId: number | null
	}>({ open: false, bookId: null })
	const [selectedCategory, setSelectedCategory] = useState<number | ''>('')

	useEffect(() => {
		refresh()
	}, [])

	const refreshData = () => refresh()

	const handleRename = async (id: number) => {
		try {
			await renameCategory(id, editName)
			setEditId(null)
			setEditName('')
			refreshData()
		} catch (error: any) {
			console.error('Error renaming category:', error)
			alert(error?.response?.data?.error || error?.message || 'Не удалось переименовать категорию. Проверьте вашу аутентификацию.')
		}
	}

	const handleCreate = async () => {
		if (!newCategory.trim()) return
		try {
			await createCategory(newCategory.trim())
			setNewCategory('')
			refreshData()
		} catch (error: any) {
			console.error('Error creating category:', error)
			alert(error?.response?.data?.error || 'Не удалось создать категорию. Проверьте вашу аутентификацию.')
		}
	}

	const handleDelete = async (id: number) => {
		if (window.confirm('Вы уверены, что хотите удалить эту категорию?')) {
			try {
				await deleteCategory(id)
				refreshData()
			} catch (error: any) {
				console.error('Error deleting category:', error)
				alert(error?.response?.data?.error || error?.message || 'Не удалось удалить категорию. Проверьте вашу аутентификацию.')
			}
		}
	}

	const handleOpenAssign = (bookId: number) => {
		setAssignDialog({ open: true, bookId })
		setSelectedCategory('')
	}

	const handleAssign = async () => {
		if (assignDialog.bookId && selectedCategory) {
			try {
				await assignBookToCategory(assignDialog.bookId, selectedCategory)
				setAssignDialog({ open: false, bookId: null })
				setSelectedCategory('')
				refreshData()
			} catch (error: any) {
				console.error('Error assigning book to category:', error)
				alert(error?.response?.data?.error || error?.message || 'Не удалось назначить книгу категории. Проверьте вашу аутентификацию.')
			}
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
				Категории
			</Typography>
			<Stack 
				direction={{ xs: 'column', sm: 'row' }} 
				spacing={2} 
				mb={4} 
				alignItems='center'
				sx={{
					p: 3,
					borderRadius: 3,
					backgroundColor: 'background.paper',
					boxShadow: '0 4px 20px rgba(2, 6, 23, 0.08)',
				}}
			>
				<TextField
					label='Новая Категория'
					value={newCategory}
					onChange={e => setNewCategory(e.target.value)}
					fullWidth
					sx={{
						'& .MuiOutlinedInput-root': {
							borderRadius: 2,
						},
					}}
				/>
				<Button 
					variant='contained' 
					onClick={handleCreate}
					sx={{
						py: 1.5,
						px: 4,
						borderRadius: 2,
						background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
						'&:hover': {
							background: 'linear-gradient(135deg, #1565c0 0%, #6a1b9a 100%)',
						},
					}}
				>
					Добавить Категорию
				</Button>
			</Stack>
			{(categories || [])
	.filter((cat: any) => !cat.is_system)
	.map((cat: any) => (

				<Box
					key={cat.id}
					sx={{ 
						mb: 4, 
						p: 3, 
						border: '2px solid rgba(25, 118, 210, 0.1)', 
						borderRadius: 3,
						backgroundColor: 'background.paper',
						boxShadow: '0 4px 20px rgba(2, 6, 23, 0.08)',
						transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
						'&:hover': {
							boxShadow: '0 8px 30px rgba(25, 118, 210, 0.12)',
							borderColor: 'rgba(25, 118, 210, 0.2)',
						},
					}}
				>
					<Stack direction='row' alignItems='center' spacing={2}>
						{editId === cat.id ? (
						<>
							<TextField
								value={editName}
								onChange={e => setEditName(e.target.value)}
							/>
							<Button 
								onClick={() => handleRename(cat.id)}
								variant='contained'
								sx={{
									borderRadius: 2,
									px: 3,
									background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
									'&:hover': {
										background: 'linear-gradient(135deg, #1565c0 0%, #6a1b9a 100%)',
									},
								}}
							>
								Сохранить
							</Button>
							<Button 
								onClick={() => setEditId(null)}
								variant='outlined'
								sx={{
									borderRadius: 2,
									px: 3,
									borderWidth: 2,
									'&:hover': {
										borderWidth: 2,
									},
								}}
							>
								Отмена
							</Button>

							{!cat.is_system && (
								<Button color='error' onClick={() => handleDelete(cat.id)}>
									Удалить
								</Button>
							)}
						</>
					) : (
							<>
								<Typography variant='h5'>{cat.name}</Typography>

								{!cat.is_system && (
									<>
										<Button
											onClick={() => {
												setEditId(cat.id)
												setEditName(cat.name)
											}}
										>
											Переименовать
										</Button>
										<Button color='error' onClick={() => handleDelete(cat.id)}>
											Удалить
										</Button>
									</>
								)}
							</>

						)}
					</Stack>
					<Box display='flex' flexWrap='wrap' gap={2} mt={2}>
						{(cat.books || []).map((book: any) => (
							<Card
								key={book.id}
								sx={{ 
									width: 200, 
									p: 2, 
									position: 'relative',
									transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
									'&:hover': {
										transform: 'translateY(-4px)',
										boxShadow: '0 8px 24px rgba(2, 6, 23, 0.15)',
									},
								}}
							>
								<Typography variant='subtitle1'>{book.title}</Typography>
								<Typography variant='body2'>{book.author}</Typography>
								<Button
									size='small'
									variant='outlined'
									sx={{ 
										position: 'absolute', 
										top: 8, 
										right: 8,
										borderRadius: 2,
										borderWidth: 2,
										'&:hover': {
											borderWidth: 2,
										},
									}}
									onClick={() => handleOpenAssign(book.id)}
								>
									Назначить
								</Button>
							</Card>
						))}
					</Box>
				</Box>
			))}
			<Box mb={5}>
				<Typography variant='h6' gutterBottom>
					Нераспределенные книги
				</Typography>
				<Box display='flex' flexWrap='wrap' gap={2} mt={2}>
					{uncategorizedBooks.map((book: any) => (
						<Card key={book.id} sx={{ width: 200, p: 2, position: 'relative' }}>
							<Typography variant='subtitle1'>{book.title}</Typography>
							<Typography variant='body2'>{book.author}</Typography>
							<Button
								size='small'
								sx={{ position: 'absolute', top: 8, right: 8 }}
								onClick={() => handleOpenAssign(book.id)}
							>
								Назначать
							</Button>
						</Card>
					))}
				</Box>
			</Box>
			<Dialog
				open={assignDialog.open}
				onClose={() => setAssignDialog({ open: false, bookId: null })}
			>
				<DialogTitle>Назначить книгу категории</DialogTitle>
				<DialogContent>
					<FormControl fullWidth>
						<InputLabel>Категории</InputLabel>
						<Select
							value={selectedCategory}
							label='Категория'
							onChange={e => setSelectedCategory(Number(e.target.value))}
						>
							{(categories || []).map(cat => (
								<MenuItem key={cat.id} value={cat.id}>
									{cat.name}
								</MenuItem>
							))}
						</Select>
					</FormControl>
				</DialogContent>
				<DialogActions>
					<Button
						onClick={() => setAssignDialog({ open: false, bookId: null })}
					>
						Отмена
					</Button>
					<Button 
						onClick={handleAssign} 
						variant='contained'
						sx={{
							borderRadius: 2,
							px: 3,
							background: 'linear-gradient(135deg, #1976d2 0%, #7c4dff 100%)',
							'&:hover': {
								background: 'linear-gradient(135deg, #1565c0 0%, #6a1b9a 100%)',
							},
						}}
					>
						Назначить
					</Button>
				</DialogActions>
			</Dialog>
		</Box>
	)
}

export default CategoriesPanel
