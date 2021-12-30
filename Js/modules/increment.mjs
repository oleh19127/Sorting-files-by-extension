// import default package
import path from 'path'
import fs from 'fs'

const increment = (dir, name, ext) => {
  let newPath = path.join(dir, `${name}${ext}`)
  let num = 1
  while (fs.existsSync(newPath)) {
    newPath = path.join(dir, `${name}(${num++})${ext}`)
    if (!fs.existsSync(newPath)) {
      return newPath
    }
  }
  return newPath
}

export { increment }