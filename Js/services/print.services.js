import chalk from 'chalk'

const log = console.log

class Print {
  successful(message) {
    log(chalk.green(message))
  }
  warning(message) {
    log(chalk.yellow(message))
  }
  error(message) {
    log(chalk.red(message))
  }
}

const print = new Print

export { print }

