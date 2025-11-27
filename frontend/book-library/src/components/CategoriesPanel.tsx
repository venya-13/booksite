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
		await renameCategory(id, editName)
		setEditId(null)
		setEditName('')
		refreshData()
	}

	const handleCreate = async () => {
		if (!newCategory.trim()) return
		await createCategory(newCategory.trim())
		setNewCategory('')
		refreshData()
	}

	const handleDelete = async (id: number) => {
		if (window.confirm('Are you sure you want to delete this category?')) {
			await deleteCategory(id)
			refreshData()
		}
	}

	const handleOpenAssign = (bookId: number) => {
		setAssignDialog({ open: true, bookId })
		setSelectedCategory('')
	}

	const handleAssign = async () => {
		if (assignDialog.bookId && selectedCategory) {
			await assignBookToCategory(assignDialog.bookId, selectedCategory)
			setAssignDialog({ open: false, bookId: null })
			setSelectedCategory('')
			refreshData()
		}
	}

	return (
		<Box>
			<Typography variant='h4' gutterBottom>
				Категории
			</Typography>
			<Stack direction='row' spacing={2} mb={3} alignItems='center'>
				<TextField
					label='Новая Категория'
					value={newCategory}
					onChange={e => setNewCategory(e.target.value)}
				/>
				<Button variant='contained' onClick={handleCreate}>
					Добавить Категорию
				</Button>
			</Stack>
			{(categories || []).map((cat: any) => (
				<Box
					key={cat.id}
					sx={{ mb: 4, p: 2, border: '1px solid #eee', borderRadius: 2 }}
				>
					<Stack direction='row' alignItems='center' spacing={2}>
						{editId === cat.id ? (
							<>
								<TextField
									value={editName}
									onChange={e => setEditName(e.target.value)}
								/>
								<Button onClick={() => handleRename(cat.id)}>Save</Button>
								<Button onClick={() => setEditId(null)}>Cancel</Button>
								<Button color='error' onClick={() => handleDelete(cat.id)}>
									Удалить
								</Button>
							</>
						) : (
							<>
								<Typography variant='h5'>{cat.name}</Typography>
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
					</Stack>
					<Box display='flex' flexWrap='wrap' gap={2} mt={2}>
						{(cat.books || []).map((book: any) => (
							<Card
								key={book.id}
								sx={{ width: 200, p: 2, position: 'relative' }}
							>
								<Typography variant='subtitle1'>{book.title}</Typography>
								<Typography variant='body2'>{book.author}</Typography>
								<Button
									size='small'
									sx={{ position: 'absolute', top: 8, right: 8 }}
									onClick={() => handleOpenAssign(book.id)}
								>
									Assign
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
				<DialogTitle>Assign Book to Category</DialogTitle>
				<DialogContent>
					<FormControl fullWidth>
						<InputLabel>Категории</InputLabel>
						<Select
							value={selectedCategory}
							label='Category'
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
					<Button onClick={handleAssign} variant='contained'>
						Назначить
					</Button>
				</DialogActions>
			</Dialog>
		</Box>
	)
}

export default CategoriesPanel
