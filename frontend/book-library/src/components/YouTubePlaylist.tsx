import React, { useState } from 'react'
import {
	Box,
	Typography,
	List,
	ListItemButton,
	ListItemText,
	Divider,
} from '@mui/material'
import { motion } from 'framer-motion'

const sections = [
	{
		title: 'Основной плейлист',
		videos: [
			{ id: 'iUvEM-USYS4', title: 'Новый Завет – лучший завет' },
			{ id: 'lnsndlan9kU', title: 'Мессианик 101' },
			{ id: 'tsznV5rQq-w', title: 'Деяние Святых' },
			{ id: 'xl1c3EPkQqU', title: 'Преданность и Вера' },
			{ id: 'vs-pQZoo-ic', title: 'С Рождеством?' },
			{ id: 'EnU_-38FBuA', title: 'Ханука: Сокрытая Истина' },
			{ id: '7OKtqIOyZUs', title: 'Кто Есть Израиль?' },
		],
	},
	{
		title: 'Дополнительные видео',
		videos: [
			{ id: 'sXUyohNQeWI', title: 'Всеоружие Божие' },
			{ id: '1MqLFkanjeE', title: 'ловушки в мессианской вере' },
		],
	},
]

const YouTubePlaylist: React.FC = () => {
	const [activeVideos, setActiveVideos] = useState(sections.map(() => 0))

	const handleVideoSelect = (sectionIdx: number, videoIdx: number) => {
		setActiveVideos(prev => {
			const updated = [...prev]
			updated[sectionIdx] = videoIdx
			return updated
		})
	}

	return (
		<Box
			sx={{
				py: 6,
				px: 2,
				background: 'linear-gradient(135deg, #f9fafb 0%, #ffffff 100%)',
				minHeight: '100vh',
			}}
		>
			<Typography variant='h3' align='center' sx={{ fontWeight: 800, mb: 6 }}>
				🎥 Видео тематики
			</Typography>

			{sections.map((sec, secIdx) => {
				const currentVideo = sec.videos[activeVideos[secIdx]]
				return (
					<Box key={secIdx} sx={{ mb: 8 }}>
						<Typography
							variant='h4'
							sx={{
								fontWeight: 700,
								mb: 2,
								background: 'linear-gradient(to right, #4f46e5, #7c3aed)',
								WebkitBackgroundClip: 'text',
								color: 'transparent',
							}}
						>
							{sec.title}
						</Typography>

						<Box
							sx={{
								display: 'flex',
								flexDirection: { xs: 'column', md: 'row' },
								gap: 3,
							}}
						>
							{/* video */}
							<motion.div
								key={currentVideo.id}
								initial={{ opacity: 0, scale: 0.96 }}
								animate={{ opacity: 1, scale: 1 }}
								transition={{ duration: 0.4 }}
								style={{
									flex: 3,
									display: 'flex',
									flexDirection: 'column',
									alignItems: 'center',
								}}
							>
								<Box
									sx={{
										width: '100%',
										maxWidth: 900,
										aspectRatio: '16/9',
										boxShadow: 4,
										borderRadius: 2,
										overflow: 'hidden',
										transition: 'box-shadow 0.3s ease',
										'&:hover': { boxShadow: 8 },
									}}
								>
									<iframe
										width='100%'
										height='100%'
										src={`https://www.youtube.com/embed/${currentVideo.id}`}
										title={currentVideo.title}
										frameBorder='0'
										allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture'
										allowFullScreen
									></iframe>
								</Box>
								<Typography
									variant='h6'
									sx={{
										mt: 2,
										fontWeight: 600,
										textAlign: 'center',
										color: '#111827',
									}}
								>
									{currentVideo.title}
								</Typography>
							</motion.div>

							{/* playlist */}
							<Box
								sx={{
									flex: 1,
									bgcolor: 'white',
									boxShadow: 3,
									borderRadius: 3,
									overflowY: 'auto',
									maxHeight: { md: '60vh', xs: 'auto' },
									minWidth: 280,
								}}
							>
								<Typography
									variant='h6'
									align='center'
									sx={{
										p: 2,
										borderBottom: '1px solid #e5e7eb',
										fontWeight: 600,
									}}
								>
									{sec.title}
								</Typography>
								<Divider />
								<List dense>
									{sec.videos.map((video, vIdx) => (
										<ListItemButton
											key={video.id}
											onClick={() => handleVideoSelect(secIdx, vIdx)}
											selected={activeVideos[secIdx] === vIdx}
											sx={{
												'&.Mui-selected': { bgcolor: '#e0e7ff' },
												'&:hover': { bgcolor: '#eef2ff' },
												transition: 'background-color 0.2s ease',
											}}
										>
											<Box
												component='img'
												src={`https://img.youtube.com/vi/${video.id}/default.jpg`}
												alt={video.title}
												sx={{
													width: 100,
													height: 60,
													objectFit: 'cover',
													borderRadius: 1,
													mr: 2,
												}}
											/>
											<ListItemText
												primary={video.title}
												primaryTypographyProps={{
													fontSize: 14,
													fontWeight: 500,
												}}
											/>
										</ListItemButton>
									))}
								</List>
							</Box>
						</Box>
					</Box>
				)
			})}
		</Box>
	)
}

export default YouTubePlaylist
