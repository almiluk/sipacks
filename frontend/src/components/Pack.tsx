import { Pack } from '../models/pack';

import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import { CardHeader, Chip, Button } from '@mui/material';
import DownloadRoundedIcon from '@mui/icons-material/DownloadRounded';

interface IPackComponentProps {
	pack: Pack,
}

export function PackComponent({ pack }: IPackComponentProps) {
	const formatDate = (date: Date) => {
		// TODO: 2-digit day and month
		return `${date.getDate()}.${date.getMonth()}.${date.getFullYear()}`
	}

	const tags = pack.tags.map((tag, i) => {
		// TODO: add link to search by tag
		return <Chip key={i} label={tag} variant="outlined" style={styles.TagChip} />
	})


	return (
		<Card style={styles.PackCard}>
			<CardHeader
				title={
					<div style={styles.PackCardHeader}>
						<span style={styles.PackName} >{pack.name}</span>
						{/* TODO: increase local download count */}
						<Button href="www.google.com" variant="outlined" style={styles.DownloadButton}>
							<span>{pack.downloadsNum}</span>
							<DownloadRoundedIcon />
							<span>{Math.floor(pack.fileSize / 1024 / 1024)} MB</span>
						</Button>
					</div>
				}
				subheader={
					<div style={styles.PackCardSubheader}>
						{`${pack.author} \u2022 ${formatDate(pack.creationDate)}`}
					</div>
				}
			/>

			<CardContent>
				<div>{tags}</div>
			</CardContent>
		</Card>
	)
}

const styles = {
	PackCard: {
		width: "100%",
	},
	DownloadButton: {
		fontSize: "large",
	},
	PackName: {
		//alignSelf: "center",
	},
	PackCardHeader: {
		display: "flex",
		justifyContent: "space-between",
	},
	PackCardSubheader: {
		textAlign: "center" as const,
	},
	TagChip: {
		margin: "0.1em",
	}
}
