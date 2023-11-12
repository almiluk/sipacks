import { Pack } from "../models/pack"
import { PackComponent } from "./Pack"
import { List, ListItem } from "@mui/material"
import { BackendRequester } from "../api/sipacks"
import { Interface } from "readline"

interface IPackListProps {
	packs: Pack[],
}

export function PackListComponent({ packs }: IPackListProps) {
	/* let packs: Pack[] = [
		{
			name: "test",
			author: "author",
			creationDate: new Date(),
			tags: ["tag1", "tag2", "long tag name", "tag3", "tag4", "tag6", "tag7", "tagtagtagtagtagtagtagtagtag", "tag8", "tag9", "tag10"],
			downloads_num: 9,
			file_size: 10.67 * 1024 * 1024,
			guid: "guid",
		},
		{
			name: "test2",
			author: "author2",
			creationDate: new Date(),
			tags: ["tag1", "tag2"],
			downloads_num: 98,
			file_size: 0,
			guid: "guid2",
		},
	] */

	let packComponents = packs.map((pack, i) => {
		return <ListItem style={styles.PackListItem} >
			<PackComponent
				key={i}
				pack={pack}
			/>
		</ListItem>
	})

	return (
		<List style={styles.PackList}>
			{packComponents}
		</List>
	);
}


const styles = {
	PackList: {
		maxWidth: "50%",
		margin: "auto",
	},
	PackListItem: {
		width: "100%",
		justifyContent: "center",
	}
}
