// import npm package
import globPkg from 'glob'
const { glob } = globPkg

// import module
import { increment } from './increment.mjs'

// import default package
import path from 'path'
import { log } from 'console'
import fs from 'fs'

const sortByExtension = (obj, sortedFolder, otherFolder) => {
  return new Promise((resolve, reject) => {
    glob('**/*', { 'ignore': ['node_modules/**/*', 'modules/**/*', 'main.js', 'modules', 'node_modules', 'package-lock.json', 'package.json', `${sortedFolder}/**/*`, `${sortedFolder}`, `${otherFolder}/**/*`, `${otherFolder}`, `${sortedFolder}.zip`, `${otherFolder}.zip`] }, (err, allFiles) => {
      if (err) {
        log(err)
      }
      if (!obj) {
        reject('No template to sort')
      }
      for (const key in obj) {
        if (Object.hasOwnProperty.call(obj, key)) {
          const data = obj[key];
          allFiles.forEach(file => {
            const ext = path.extname(file)
            const name = path.basename(file, ext)
            const extName = ext.replace('.', '')
            data.extensions.forEach(extension => {
              if (extName.toLowerCase() === extension.toLowerCase()) {
                const statFile = fs.statSync(file)
                const dataCreatedFile = statFile.mtime.getFullYear().toString()
                if (!fs.existsSync(path.join(sortedFolder, dataCreatedFile, data.folder))) {
                  fs.mkdirSync(path.join(sortedFolder, dataCreatedFile, data.folder), { recursive: true })
                }
                const newPath = increment(path.join(sortedFolder, dataCreatedFile, data.folder), name, ext)
                fs.renameSync(file, newPath)
                log(`${file} moved to ${newPath}`)
              }
            })
          })
        }
      }
      resolve('')
    })
  })
}

export { sortByExtension }