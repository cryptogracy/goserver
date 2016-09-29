package main

func stdErrCheck(err error) {
  if err != nil {
    panic(err)
  }
}
