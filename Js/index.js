import { sort } from "./services/sorting.services.js"

const init = async () => {
  await sort.byExtension()
  await sort.removeAllEmptyFolders()
}

init()
