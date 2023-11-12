import { AppBar, TextField, IconButton, } from '@mui/material'
import SearchIcon from "@mui/icons-material/Search"

import { Pack } from "../models/pack"

interface ITopBarProps {
	findPacks: () => void,
}

export function TopBar({ findPacks }: ITopBarProps) {
	return (
		<AppBar position="static">
			<div style={styles.TopBarContainer}>
				<h1>SIPacks</h1>
				<form style={{ display: "flex", alignItems: "center" }}>
					<TextField
						sx={styles.SearchBarInput}
						placeholder="Pack name..."
						InputProps={{
							endAdornment: (
								<IconButton style={styles.SearchBarIcon} onClick={findPacks}>
									<SearchIcon />
								</IconButton>
							),
						}}
					/>
				</form>
			</div>
		</AppBar>
	);
}

const styles = {
	TopBarContainer: {
		display: "flex",
		flexDirection: "row",
		justifyContent: "space-around",
		alignItems: "stretch",
	} as const,
	SearchBar: {
		display: "flex",
		flexDirection: "row",
		justifyContent: "flex-end",
		alignItems: "center",
	},
	SearchBarInput: {
		"& .MuiOutlinedInput-root": {
			"&.Mui-focused fieldset": {
				border: "1px solid black"
			}
		}
	},
	SearchBarIcon: {
	},
}
