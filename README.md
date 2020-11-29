# ZKits Runner Library #

[![ZKits](https://img.shields.io/badge/ZKits-Library-f3c)](https://github.com/edoger/zkits-runner)
[![Build Status](https://travis-ci.org/edoger/zkits-runner.svg?branch=master)](https://travis-ci.org/edoger/zkits-runner)
[![Build status](https://ci.appveyor.com/api/projects/status/akl62co7bn4wtgvf/branch/master?svg=true)](https://ci.appveyor.com/project/edoger56924/zkits-runner/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-runner/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-runner?branch=master)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b6cfc08a46a04e19acfbf722b013567e)](https://www.codacy.com/manual/edoger/zkits-runner/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-runner&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/edoger/zkits-runner)](https://goreportcard.com/report/github.com/edoger/zkits-runner)
[![Golang Version](https://img.shields.io/badge/golang-1.13+-orange)](https://github.com/edoger/zkits-runner)

## About ##

This package is a library of ZKits project. 
This library provides a convenient subtask runner for applications. 
We can easily control the running order of subtasks and exit them in reverse order.

## Usage ##

 1. Import package.
 
    ```sh
    go get -u -v github.com/edoger/zkits-runner
    ```

 2. Create a runner to run subtasks within the application.

    ```go
    package main
    
    import (
       "github.com/edoger/zkits-runner"
    )
    
    func main() {
       r := runner.New()
       err := r.Run(runner.NewTaskFromFunc(nil, func() error {
           // Do something.
           return nil
       }))
       // Wait system exit.
       if err := r.Wait(); err != nil {
           // Handle error.
       }
    }
    ```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
