import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import App from './App'
import reportWebVitals from './reportWebVitals'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import { BrowserRouter as Router } from 'react-router-dom'

const theme = createTheme({
	palette: {
		mode: 'light',
		primary: { main: '#1976d2' },
		secondary: { main: '#7c4dff' },
		background: { default: '#f7f9fc', paper: '#ffffff' },
	},
	typography: {
		fontFamily:
			'Inter, system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif',
		h1: { fontWeight: 700, letterSpacing: '-0.02em' },
		h2: { fontWeight: 700, letterSpacing: '-0.01em' },
		h3: { fontWeight: 700, letterSpacing: '-0.01em' },
		h4: { fontWeight: 700, letterSpacing: '-0.005em' },
		h5: { fontWeight: 600, letterSpacing: '0em' },
		h6: { fontWeight: 600, letterSpacing: '0.01em' },
		button: { textTransform: 'none', fontWeight: 600, letterSpacing: '0.02em' },
		body1: { lineHeight: 1.7 },
		body2: { lineHeight: 1.6 },
	},
	shape: {
		borderRadius: 16,
	},
	components: {
		MuiButton: {
			styleOverrides: {
				root: {
					paddingLeft: 20,
					paddingRight: 20,
					paddingTop: 10,
					paddingBottom: 10,
					boxShadow: '0 2px 8px rgba(25, 118, 210, 0.15)',
					transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
					'&:hover': {
						boxShadow: '0 4px 16px rgba(25, 118, 210, 0.25)',
						transform: 'translateY(-2px)',
					},
					'&:active': {
						transform: 'translateY(0px)',
					},
				},
			},
		},
		MuiPaper: {
			styleOverrides: {
				root: {
					borderRadius: 20,
					boxShadow: '0 4px 20px rgba(2, 6, 23, 0.08)',
					transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
				},
			},
		},
		MuiCard: {
			styleOverrides: {
				root: {
					borderRadius: 20,
					boxShadow: '0 4px 20px rgba(2, 6, 23, 0.08)',
					transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
					overflow: 'hidden',
					'&:hover': {
						transform: 'translateY(-6px)',
						boxShadow: '0 12px 40px rgba(2, 6, 23, 0.15)',
					},
				},
			},
		},
		MuiTextField: {
			styleOverrides: {
				root: {
					'& .MuiOutlinedInput-root': {
						transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
						'&:hover': {
							transform: 'translateY(-1px)',
							boxShadow: '0 4px 12px rgba(25, 118, 210, 0.1)',
						},
						'&.Mui-focused': {
							transform: 'translateY(-1px)',
							boxShadow: '0 4px 16px rgba(25, 118, 210, 0.2)',
						},
					},
				},
			},
		},
	},
})

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement)
root.render(
	<React.StrictMode>
		<ThemeProvider theme={theme}>
			<CssBaseline />
			<Router>
				<App />
			</Router>
		</ThemeProvider>
	</React.StrictMode>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
