export class Pack {
	name: string
	author: string
	creationDate: Date
	guid: string
	fileSize: number
	downloadsNum: number
	tags: string[]

	constructor(name: string, author: string, creation_date: Date, guid: string, file_size: number, downloads_num: number, tags: string[]) {
		this.name = name
		this.author = author
		this.creationDate = creation_date
		this.guid = guid
		this.fileSize = file_size
		this.downloadsNum = downloads_num
		this.tags = tags
	}
}
