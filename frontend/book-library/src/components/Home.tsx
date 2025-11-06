import { motion } from 'framer-motion'
import { Button } from '@mui/material'

export default function Home() {
	return (
		<div
			style={{
				minHeight: '80vh',
				display: 'flex',
				flexDirection: 'column',
				alignItems: 'center',
				justifyContent: 'center',
				textAlign: 'center',
				background: 'linear-gradient(180deg, #f5f7ff 0%, #eaf1ff 100%)',
				padding: '4rem 2rem',
			}}
		>
			<motion.div
				initial={{ opacity: 0, y: 30 }}
				animate={{ opacity: 1, y: 0 }}
				transition={{ duration: 1 }}
				style={{ maxWidth: '1000px', width: '100%' }}
			>
				<motion.h1
					initial={{ opacity: 0, y: 20 }}
					animate={{ opacity: 1, y: 0 }}
					transition={{ delay: 0.3, duration: 1 }}
					style={{
						fontSize: '3.5rem',
						fontWeight: 900,
						background: 'linear-gradient(90deg, #1976d2, #7c4dff)',
						WebkitBackgroundClip: 'text',
						WebkitTextFillColor: 'transparent',
						marginBottom: '1.5rem',
						lineHeight: 1.2,
					}}
				>
					ЦЕНТР УЧЕНИЧЕСТВА
					<br />
					СЛУЖЕНИЕ АРТУРА БЭЙЛИ
				</motion.h1>

				<motion.p
					initial={{ opacity: 0 }}
					animate={{ opacity: 1 }}
					transition={{ delay: 0.8, duration: 1 }}
					style={{
						fontSize: '1.4rem',
						color: '#2b2b2b',
						maxWidth: '800px',
						margin: '0 auto 2.5rem',
						lineHeight: 1.6,
					}}
				>
					Наша цель — донести Истинное Евангелие до краёв земли и восстановить
					веру, однажды преданную святым.
				</motion.p>

				<motion.div
					initial={{ opacity: 0 }}
					animate={{ opacity: 1 }}
					transition={{ delay: 1.4, duration: 1 }}
					style={{
						display: 'flex',
						gap: '1.5rem',
						justifyContent: 'center',
						flexWrap: 'wrap',
					}}
				>
					<Button
						variant='contained'
						size='large'
						href='/youtube'
						sx={{
							background: 'linear-gradient(90deg, #1976d2, #7c4dff)',
							fontWeight: 700,
							fontSize: '1.1rem',
							px: 4,
							py: 1.5,
							borderRadius: '9999px',
							transition: 'all .3s ease',
							'&:hover': {
								transform: 'translateY(-3px)',
								background: 'linear-gradient(90deg, #1565c0, #6a1b9a)',
							},
						}}
					>
						Смотреть видео
					</Button>

					<Button
						variant='outlined'
						size='large'
						href='/about'
						sx={{
							borderColor: '#7c4dff',
							color: '#7c4dff',
							fontWeight: 700,
							fontSize: '1.1rem',
							px: 4,
							py: 1.5,
							borderRadius: '9999px',
							transition: 'all .3s ease',
							'&:hover': {
								transform: 'translateY(-3px)',
								borderColor: '#1976d2',
								color: '#1976d2',
							},
						}}
					>
						Узнать больше
					</Button>

					<Button
						variant='outlined'
						size='large'
						href='/books'
						sx={{
							borderColor: '#1976d2',
							color: '#1976d2',
							fontWeight: 700,
							fontSize: '1.1rem',
							px: 4,
							py: 1.5,
							borderRadius: '9999px',
							transition: 'all .3s ease',
							'&:hover': {
								transform: 'translateY(-3px)',
								borderColor: '#7c4dff',
								color: '#7c4dff',
							},
						}}
					>
						Перейти к книгам
					</Button>
				</motion.div>
			</motion.div>
		</div>
	)
}
