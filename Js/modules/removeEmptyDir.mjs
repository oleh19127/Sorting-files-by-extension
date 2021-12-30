// import default package
import { log } from 'console'
import fs from 'fs'

const removeEmptyDir = (dirs) => {
  for (let i = 0; i < dirs.length; i++) {
    dirs.forEach(dir => {
      if (fs.existsSync(dir)) {
        const filesInDir = fs.readdirSync(dir)
        if (filesInDir.length === 0) {
          fs.rmSync(dir, { recursive: true })
          log(`Removed ${dir} directory`)
        }
      }
    })
  }
}

export { removeEmptyDir }