// import modules
import { allData } from './modules/structureOfExtension.mjs'
import { sortOtherFiles } from './modules/sortingOtherFiles.mjs'
import { sortByExtension } from './modules/sortingByExtension.mjs'

// import default package
import { log } from 'console'

// const
const sortedFilesFolder = 'Sorted Files'
const otherFilesFolder = 'Other Files'

sortByExtension(allData, sortedFilesFolder, otherFilesFolder)
  .then(data => {
    log(data)
  })
  .then(() => {
    sortOtherFiles(sortedFilesFolder, otherFilesFolder)
  })
