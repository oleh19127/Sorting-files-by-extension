// import npm package
import globPkg from 'glob'
const { glob } = globPkg

// import module
import { increment } from './increment.mjs'
import { removeEmptyDir } from './removeEmptyDir.mjs'

// import default package
import path from 'path'
import { log } from 'console'
import fs from 'fs'

const sortOtherFiles = (sortedFolder, otherFolder) => {
  glob('**/*', { 'ignore': ['node_modules/**/*', 'modules/**/*', 'main.js', 'modules', 'node_modules', 'package-lock.json', 'package.json', `${sortedFolder}/**/*`, `${sortedFolder}`, `${otherFolder}/**/*`, `${otherFolder}`, `${sortedFolder}.zip`, `${otherFolder}.zip`] }, (err, allFiles) => {
    let dirs = []
    if (err) {
      log(err)
    }
    allFiles.forEach(file => {
      const statFile = fs.statSync(file)
      if (!statFile.isDirectory()) {
        const ext = path.extname(file)
        const name = path.basename(file, ext)
        const extName = ext.replace('.', '')
        if (!fs.existsSync(path.join(otherFolder, extName))) {
          fs.mkdirSync(path.join(otherFolder, extName), { recursive: true })
        }
        const newPath = increment(path.join(otherFolder, extName), name, ext)
        fs.renameSync(file, newPath)
        log(`${file} moved to ${newPath}`)
      } else if (statFile.isDirectory()) {
        dirs.push(file)
      }
    })
    removeEmptyDir(dirs)
  })
}

export { sortOtherFiles }