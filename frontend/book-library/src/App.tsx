import React, { useState } from 'react'
import {
	BrowserRouter as Router,
	Routes,
	Route,
	Link,
	useNavigate,
} from 'react-router-dom'
import {
	AppBar,
	Toolbar,
	Typography,
	Button,
	Container,
	Box,
	IconButton,
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

function App() {
	const [user, setUser] = useState<string | null>(null)
	const [isAdmin, setIsAdmin] = useState(false)
	const [error, setError] = useState<string | null>(null)

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
	}

	function Login() {
		const [loading, setLoading] = useState(false)
		const [err, setErr] = useState<string | null>(null)
		const navigate = useNavigate()
		const handleLogin = async (email: string, password: string) => {
			setLoading(true)
			setErr(null)
			try {
				const res = await login(email, password)
				setToken(res.token)
				const decoded: any = jwtDecode(res.token)
				setUser(decoded.email)
				setIsAdmin(!!decoded.is_admin)
				navigate('/')
			} catch (e: any) {
				setErr(e.response?.data?.error || 'Ошибка входа')
			} finally {
				setLoading(false)
			}
		}
		const handleGoogleAuth = () => {
			window.location.href = 'http://localhost:8080/login'
		}

		return (
			<Box p={2}>
				<AuthForm
					type='login'
					onSubmit={handleLogin}
					onGoogleAuth={handleGoogleAuth}
					onlyGoogle
				/>
				{err && <Typography color='error'>{err}</Typography>}
				{loading && <Typography>Загрузка...</Typography>}
			</Box>
		)
	}
	function Register() {
		const [loading, setLoading] = useState(false)
		const [err, setErr] = useState<string | null>(null)
		const navigate = useNavigate()
		const handleRegister = async (email: string, password: string) => {
			setLoading(true)
			setErr(null)
			try {
				const res = await register(email, password)
				setToken(res.token)
				const decoded: any = jwtDecode(res.token)
				setUser(decoded.email)
				setIsAdmin(!!decoded.is_admin)
				navigate('/')
			} catch (e: any) {
				setErr(e.response?.data?.error || 'Регистрация не удалась')
			} finally {
				setLoading(false)
			}
		}
		const handleGoogleAuth = () => {
			alert('Google Sign-In (mock)')
		}
		return (
			<Box p={2}>
				<AuthForm
					type='register'
					onSubmit={handleRegister}
					onGoogleAuth={handleGoogleAuth}
				/>
				{err && <Typography color='error'>{err}</Typography>}
				{loading && <Typography>Загрузка...</Typography>}
			</Box>
		)
	}

	return (
		<Router>
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
						sx={{ mr: 1 }}
					>
						<MenuBookIcon />
					</IconButton>
					<Typography
						variant='h6'
						sx={{ flexGrow: 1, fontWeight: 800, letterSpacing: 0.2 }}
					>
						Book Library
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
						Про нас
					</Button>
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
							Admin
						</Button>
					)}
				</Toolbar>
			</AppBar>
			<Container maxWidth='lg' sx={{ mt: 4, mb: 8 }}>
				<Routes>
					<Route path='/' element={<Home />} />
					<Route path='/books' element={<Books />} />
					<Route path='/about' element={<About />} />
					<Route path='/youtube' element={<YouTubePlaylist />} />
					<Route path='/login' element={<Login />} />
					<Route path='/register' element={<Register />} />
					<Route path='/admin' element={<AdminPanel />} />
					<Route path='/categories' element={<CategoriesPanel />} />
				</Routes>
			</Container>
		</Router>
	)
}

export default App
