export type RPCMethods =
  | "Service.Exec"
  | "Service.Kill"
  | "Service.Clear"
  | "Service.Running"
  | "Service.KillAll"
  | "Service.FreeSpace"
  | "Service.Formats"
  | "Service.ExecPlaylist"
  | "Service.DirectoryTree"
  | "Service.UpdateExecutable"
  | "Service.ExecLivestream"
  | "Service.ProgressLivestream"
  | "Service.KillLivestream"
  | "Service.KillAllLivestream"
  | "Service.ClearCompleted"

export type RPCRequest = {
  method: RPCMethods
  params?: any[]
  id?: string
}

export type RPCResponse<T> = Readonly<{
  result: T
  error: number | null
  id?: string
}>

type DownloadInfo = {
  url: string
  filesize_approx?: number
  resolution?: string
  thumbnail: string
  title: string
  vcodec?: string
  acodec?: string
  ext?: string
  created_at: string
}

export enum ProcessStatus {
  PENDING = 0,
  DOWNLOADING,
  COMPLETED,
  ERRORED,
  LIVESTREAM,
}

type DownloadProgress = {
  speed: number
  eta: number
  percentage: string
  process_status: ProcessStatus
}

export type RPCResult = Readonly<{
  id: string
  progress: DownloadProgress
  info: DownloadInfo
  output: {
    savedFilePath: string
  }
}>

export type RPCParams = {
  URL: string
  Params?: string
}

export type DLMetadata = {
  formats: Array<DLFormat>
  _type: string
  best: DLFormat
  thumbnail: string
  title: string
  entries: Array<DLMetadata>
}

export type DLFormat = {
  format_id: string
  format_note: string
  fps: number
  resolution: string
  vcodec: string
  acodec: string
  filesize_approx: number
  language: string
}

export type DirectoryEntry = {
  name: string
  path: string
  size: number
  modTime: string
  isVideo: boolean
  isDirectory: boolean
}

export type DeleteRequest = Pick<DirectoryEntry, 'path'>

export type PlayRequest = DeleteRequest

export type CustomTemplate = {
  id: string
  name: string
  content: string
}

export enum LiveStreamStatus {
  WAITING,
  IN_PROGRESS,
  COMPLETED,
  ERRORED
}

export type LiveStreamProgress = Record<string, {
  status: LiveStreamStatus
  waitTime: string
  liveDate: string
}>

export type RPCVersion = {
  rpcVersion: string
  ytdlpVersion: string
}

export type ArchiveEntry = {
  id: string
  title: string
  path: string
  thumbnail: string
  source: string
  metadata: string
  created_at: string
}

export type PaginatedResponse<T> = {
  first: number
  next: number
  data: T
}