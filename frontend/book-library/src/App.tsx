import React, { useState } from 'react'
import {
	AppBar,
	Toolbar,
	Typography,
	Button,
	Container,
	Box,
	IconButton,
	TextField,
} from '@mui/material'
import MenuBookIcon from '@mui/icons-material/MenuBook'
import Books from './components/Books'
import AuthForm from './components/AuthForm'
import AdminPanel from './components/AdminPanel'
import CategoriesPanel from './components/CategoriesPanel'
import { register, login, setToken, removeToken, getToken } from './api'
// Try both import styles for compatibility
import { jwtDecode } from 'jwt-decode'
import YouTubePlaylist from './components/YouTubePlaylist'
import Home from './components/Home'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import { Routes, Route, Link, useNavigate, Navigate } from 'react-router-dom'
import LeadershipCourse from './components/LeadershipCourse'

function About() {
	return (
		<Box p={2}>
			<Typography
				variant='h3'
				align='center'
				gutterBottom
				sx={{ fontWeight: 'bold', color: '#1976d2' }}
			>
				О Нас
			</Typography>
			<Box
				sx={{
					background: 'linear-gradient(135deg, #e3f2fd 0%, #fff 100%)',
					borderRadius: 3,
					boxShadow: 3,
					p: 4,
					maxWidth: 800,
					mx: 'auto',
					mt: 3,
				}}
			>
				<Typography
					variant='h6'
					align='justify'
					sx={{ fontSize: '1.2rem', color: '#333', lineHeight: 1.7 }}
				>
					Наш Центр Ученичества “ЖИЗНЬ В ИСТИНЕ” – при служении Артура Бейли
					“Дом Израиля”, Мессианское Сообщество по восстановлению Ивритских
					Корней (Hebrew Roots) верования не только христиан, но и иудеев. Наша
					цель-достигать, проповедовать и учить мужчин и женщин быть учениками
					Царства Иехова. Для исполнения нашей цели мы используем метод
					интенсивных и динамических курсов, программ и учений. Мы уверены, что
					по завершению программ и учений, вы будете иметь ясное понимание
					ученичества. Понимание Слова в действии укрепит вас и позволит вам
					более эффективно изучать и распространять мощное и жизненно-важное
					Истинное Евангелие Царства Божьего Иехова
				</Typography>
			</Box>
		</Box>
	)
}
function ProgramUchenichestva() {
	const [form, setForm] = useState({ name: '', email: '', message: '' })

	const handleChange = (
		e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => {
		setForm({ ...form, [e.target.name]: e.target.value })
	}

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		alert('Сообщение отправлено!')
	}

	return (
		<Box sx={{ maxWidth: '900px', mx: 'auto', py: 4 }}>
			<Typography variant='h4' textAlign='center' mb={3}>
				Программа Ученичества
			</Typography>

			<Typography variant='body1' mb={4} sx={{ lineHeight: 1.8 }}>
				Центр Ученичества при сообществе “ДОМ ИЗРАИЛЯ” представляет
				Аккредитованную двухлетнюю “Программу Ученичество”, которая
				предназначена для “оснащения святых” и подготовки учеников для Мессии
				Иешуа. Артур Бейли является канцлером Университета, руководителем
				Международного служения Артура Бейли и апостолом, пастором, и учителем
				Дома Израиля в Шарлотте, Северная Каролина. Arthur Bailey Ministries
				International — это мессианское служение, которое обучает ивритским
				корням христианской веры.
			</Typography>
		</Box>
	)
}

function KursLiderstvo() {
	return (
		<LeadershipCourse />
	)
}

function AcademyUcheniya() {
	return (
		<Box sx={{ maxWidth: '900px', mx: 'auto', py: 4 }}>
			<Typography
				variant='h3'
				textAlign='center'
				gutterBottom
				sx={{ fontWeight: 'bold', color: '#1976d2' }}
			>
				Академ Учения
			</Typography>

			<Typography
				variant='body1'
				sx={{ lineHeight: 1.8, fontSize: '1.2rem', color: '#333' }}
			>
				Апостол Доктор Богословия Артур Бейли является разработчиком и
				продюсером ведущей программы ученичество под названием “Ученичество
				101”, и ведущей программы под названием “Лидерство 101” и Программы для
				Служителей. Др. Бейли внес ясность во многие противоречивые и трудные
				для понимания библейские отрывки. Он оснащает людей Иехова для служения,
				чтобы ответить на призыв лидерства, ученичества и служения. Исследуйте
				этот сайт, который был разработан для русскоязычной аудитории и,
				пожалуйста, воспользуйтесь многочисленными электронными книгами и
				учениями бесплатно, для вашего обучения и роста опыта в Иешуа Мессии.
			</Typography>
		</Box>
	)
}

function KursSluzhiteley() {
	const [form, setForm] = useState({ name: '', email: '', message: '' })

	const handleChange = (
		e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => {
		setForm({ ...form, [e.target.name]: e.target.value })
	}

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		alert('Сообщение отправлено!')
	}

	return (
		<Box sx={{ maxWidth: '900px', mx: 'auto', py: 4 }}>
			<Typography variant='h4' textAlign='center' mb={3}>
				Курс Подготовки Служителей
			</Typography>

			<Typography variant='body1' mb={4} sx={{ lineHeight: 1.8 }}>
				Призвание к Служению — это наивысшее Призвание, которое вы можете
				получить на земле, потому что оно исходит от Самого Творца. Вы должны
				видеть себя мировым лидером. Он призвал вас из всех людей представлять
				Его, Царя царей, Царя Славы! Важно, чтобы вы реализовали свое Призвание.
				<br />
				<br />
				В этом курсе Артур Бэйли затрагивает темы: Призвание служителя, развитие
				и подготовка; Женщины в служении; Формы служения; Дары; Великое
				поручение; Методы исследования; Личностное развитие; Границы отношений;
				Финансовая честность; Стресс.
				<br />
				<br />
				Этот курс поможет вам утвердиться в своем призвании-посланника Мессии.
			</Typography>

			<Box
				component='form'
				onSubmit={handleSubmit}
				sx={{ p: 3, border: '1px solid #ccc', borderRadius: 2 }}
			>
				<Typography variant='h6' mb={2}>
					Связаться с нами
				</Typography>

				<TextField
					fullWidth
					label='Имя'
					name='name'
					value={form.name}
					onChange={handleChange}
					margin='normal'
				/>

				<TextField
					fullWidth
					label='Email'
					name='email'
					value={form.email}
					onChange={handleChange}
					margin='normal'
				/>

				<TextField
					fullWidth
					label='Сообщение'
					name='message'
					value={form.message}
					onChange={handleChange}
					margin='normal'
					multiline
					rows={4}
				/>

				<Button type='submit' variant='contained' sx={{ mt: 2 }}>
					Отправить
				</Button>
			</Box>
		</Box>
	)
}

function App() {
	const [user, setUser] = useState<string | null>(null)
	const [anchorElAcadem, setAnchorElAcadem] = useState<null | HTMLElement>(null)
	const [isAdmin, setIsAdmin] = useState(false)
	const [error, setError] = useState<string | null>(null)
	const navigate = useNavigate()

	React.useEffect(() => {
		const params = new URLSearchParams(window.location.search)
		const tokenFromURL = params.get('token')
		if (tokenFromURL) {
			setToken(tokenFromURL)
			window.history.replaceState({}, document.title, '/')
		}
	}, [])

	React.useEffect(() => {
		const token = getToken()
		if (token) {
			try {
				const decoded: any = jwtDecode(token)
				console.log('Decoded JWT:', decoded)
				setUser(decoded.email)
				setIsAdmin(!!decoded.is_admin)
			} catch {
				setUser(null)
				setIsAdmin(false)
			}
		} else {
			setUser(null)
			setIsAdmin(false)
		}
	}, [])

	const handleLogout = () => {
		removeToken()
		setUser(null)
		setIsAdmin(false)
		navigate('/')
	}

	function Login() {
		React.useEffect(() => {
			window.location.href = 'http://localhost:8080/login'
		}, [])

		return (
			<Box p={2}>
				<Typography>Перенаправляем на Google для авторизации...</Typography>
			</Box>
		)
	}

	function Register() {
		React.useEffect(() => {
			window.location.href = 'http://localhost:8080/login'
		}, [])

		return (
			<Box p={2}>
				<Typography>Перенаправляем на Google для регистрации...</Typography>
			</Box>
		)
	}

	return (
		<>
			<AppBar
				position='sticky'
				elevation={0}
				sx={{
					background: 'linear-gradient(90deg, #1976d2 0%, #7c4dff 100%)',
					boxShadow: '0 6px 20px rgba(25,118,210,0.25)',
				}}
			>
				<Toolbar sx={{ gap: 1 }}>
					<IconButton
						edge='start'
						color='inherit'
						aria-label='logo'
						component={Link}
						to='/'
						sx={{ mr: 1 }}
					>
						<MenuBookIcon />
					</IconButton>
					<Typography
						variant='h6'
						sx={{ flexGrow: 1, fontWeight: 800, letterSpacing: 0.2 }}
					>
						Lives in truth
					</Typography>
					<Button
						color='inherit'
						component={Link}
						to='/'
						sx={{
							opacity: 0.95,
							'&:hover': { opacity: 1, transform: 'translateY(-1px)' },
							transition: 'all .15s ease',
						}}
					>
						Главная страница
					</Button>
					<Button
						color='inherit'
						component={Link}
						to='/books'
						sx={{
							opacity: 0.95,
							'&:hover': { opacity: 1, transform: 'translateY(-1px)' },
							transition: 'all .15s ease',
						}}
					>
						Книги
					</Button>
					<Button
						color='inherit'
						component={Link}
						to='/about'
						sx={{
							opacity: 0.95,
							'&:hover': { opacity: 1, transform: 'translateY(-1px)' },
							transition: 'all .15s ease',
						}}
					>
						О нас
					</Button>

					{/* Академ Учения */}
					<Button
						color='inherit'
						onClick={e => setAnchorElAcadem(e.currentTarget)}
						sx={{
							opacity: 0.95,
							'&:hover': { opacity: 1, transform: 'translateY(-1px)' },
							transition: 'all .15s ease',
						}}
					>
						Академ Учения
					</Button>

					<Menu
						anchorEl={anchorElAcadem}
						open={Boolean(anchorElAcadem)}
						onClose={() => setAnchorElAcadem(null)}
					>
						<MenuItem
							component={Link}
							to='/academy'
							onClick={() => setAnchorElAcadem(null)}
						>
							Об Академии
						</MenuItem>
						<MenuItem
							component={Link}
							to='/program-uchenichestva'
							onClick={() => setAnchorElAcadem(null)}
						>
							Программа Ученичества
						</MenuItem>
						<MenuItem
							component={Link}
							to='/kurs-liderstvo'
							onClick={() => setAnchorElAcadem(null)}
						>
							Курс Лидерство
						</MenuItem>
						<MenuItem
							component={Link}
							to='/kurs-sluzhiteley'
							onClick={() => setAnchorElAcadem(null)}
						>
							Курс Подготовки Служителей
						</MenuItem>
					</Menu>
					<Button
						color='inherit'
						component={Link}
						to='/youtube'
						sx={{
							opacity: 0.95,
							'&:hover': { opacity: 1, transform: 'translateY(-1px)' },
							transition: 'all .15s ease',
						}}
					>
						Плейлист YouTube
					</Button>
					{isAdmin && (
						<Button
							color='inherit'
							component={Link}
							to='/categories'
							sx={{
								opacity: 0.95,
								'&:hover': { opacity: 1, transform: 'translateY(-1px)' },
								transition: 'all .15s ease',
							}}
						>
							Категории
						</Button>
					)}
					{!user && (
						<Button
							color='inherit'
							component={Link}
							to='/login'
							sx={{
								ml: 1,
								borderRadius: 9999,
								px: 2,
								backgroundColor: 'rgba(255,255,255,0.12)',
								'&:hover': { backgroundColor: 'rgba(255,255,255,0.2)' },
							}}
						>
							Войти/Зарегистрироваться
						</Button>
					)}
					{user && (
						<Button
							color='inherit'
							onClick={handleLogout}
							sx={{
								ml: 1,
								borderRadius: 9999,
								px: 2,
								backgroundColor: 'rgba(255,255,255,0.12)',
								'&:hover': { backgroundColor: 'rgba(255,255,255,0.2)' },
							}}
						>
							Выход
						</Button>
					)}
					{isAdmin && (
						<Button
							color='inherit'
							component={Link}
							to='/admin'
							sx={{
								ml: 1,
								borderRadius: 9999,
								px: 2,
								backgroundColor: 'rgba(255,255,255,0.12)',
								'&:hover': { backgroundColor: 'rgba(255,255,255,0.2)' },
							}}
						>
							Админ
						</Button>
					)}
				</Toolbar>
			</AppBar>
			<Box
				sx={{
					minHeight: '100vh', // чтобы занимало всю высоту экрана
					display: 'flex',
					flexDirection: 'column',
					background: 'linear-gradient(180deg, #f5f7ff 0%, #eaf1ff 100%)', // фон на всю страницу
				}}
				>
			<Box sx={{ flex: 1, width: '100%' }}> {/* контент будет растягиваться */}
				<Routes>
				<Route path='/' element={<Home />} />
				<Route path='/books' element={<Books />} />
				<Route path='/about' element={<About />} />
				<Route path='/youtube' element={<YouTubePlaylist />} />
				<Route path='/login' element={<Login />} />
				<Route path='/register' element={<Register />} />
				<Route
					path='/admin'
					element={isAdmin ? <AdminPanel /> : <Navigate to='/' replace />}
				/>
				<Route path='/categories' element={<CategoriesPanel />} />
				<Route path='/academy' element={<AcademyUcheniya />} />
				<Route path='/program-uchenichestva' element={<ProgramUchenichestva />} />
				<Route path='/kurs-liderstvo' element={<KursLiderstvo />} />
				<Route path='/kurs-sluzhiteley' element={<KursSluzhiteley />} />
				</Routes>
			</Box>
		</Box>

		</>
	)
}

export default App

