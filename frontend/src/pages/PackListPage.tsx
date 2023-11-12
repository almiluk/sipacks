import React, { useState } from "react";

import { TopBar } from "../components/TopBar"
import { PackListComponent } from "../components/PackList"
import { BackendRequester } from "../api/sipacks"
import { Pack } from "../models/pack"


export function PackListPage() {
	const backendRequester = new BackendRequester("localhost", 8080, "api/v1")

	const [packs, setPacks] = React.useState<Pack[]>([])

	const findPacks = () => {
		backendRequester.getPacks().then((response) => {
			setPacks(response.data.packs.map((pack: any) => new Pack(
				pack.name,
				pack.author,
				new Date(pack.creation_date),
				pack.guid,
				pack.file_size,
				pack.downloads_num,
				pack.tags,
			)))
		})
	}

	return (
		<div>
			<TopBar findPacks={findPacks} />
			<PackListComponent packs={packs} />
		</div>
	);
}
