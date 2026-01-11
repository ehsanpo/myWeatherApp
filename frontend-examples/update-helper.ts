// Auto-Update Helper
import { CheckForUpdates, GetCurrentVersion, OpenReleaseURL } from '../wailsjs/go/main/App'

interface UpdateInfo {
  version: string
  releaseUrl: string
  downloadUrl: string
  description: string
  available: boolean
}

export async function checkForUpdates(): Promise<UpdateInfo | null> {
  try {
    const updateInfo = await CheckForUpdates()
    return updateInfo
  } catch (error) {
    console.error('Failed to check for updates:', error)
    return null
  }
}

export async function getCurrentVersion(): Promise<string> {
  try {
    return await GetCurrentVersion()
  } catch (error) {
    console.error('Failed to get current version:', error)
    return 'unknown'
  }
}

export async function openReleaseURL(url: string) {
  try {
    await OpenReleaseURL(url)
  } catch (error) {
    console.error('Failed to open release URL:', error)
  }
}

// Example: Check for updates and notify user
export async function checkAndNotify() {
  const updateInfo = await checkForUpdates()
  
  if (updateInfo && updateInfo.available) {
    const shouldUpdate = confirm(
      `A new version (${updateInfo.version}) is available!\n\n${updateInfo.description.substring(0, 200)}...\n\nWould you like to download it?`
    )
    
    if (shouldUpdate) {
      if (updateInfo.downloadUrl) {
        await openReleaseURL(updateInfo.downloadUrl)
      } else {
        await openReleaseURL(updateInfo.releaseUrl)
      }
    }
  } else {
    console.log('You are running the latest version')
  }
}
