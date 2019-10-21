# PROCESS PIPE

An abstraction layer around the [exec.Cmd](https://golang.org/pkg/os/exec/#Cmd) for easily creating linked commands in an [Pipeline](https://en.wikipedia.org/wiki/Pipeline_(Unix) format.

So for example:

```
cmd, _ := pipe.NewProcess(`echo -en "x\ncee\nfoo\nbar" | sort`)
   
```

will be have an signature like:

```
(*pipe.Process)((len=2 cap=2) {
 (*exec.Cmd)({
  Path: (string) (len=9) "/bin/echo",
  Args: ([]string) (len=3 cap=3) {
   (string) (len=4) "echo",
   (string) (len=3) "-en",
   (string) (len=16) "x\\ncee\\nfoo\\nbar"
  },
  .....
 }),
 (*exec.Cmd)({
  Path: (string) (len=13) "/usr/bin/sort",
  Args: ([]string) (len=1 cap=1) {
   (string) (len=4) "sort"
  },
 .....
 })
})

```

where the stdout is linked to next commands stdin. 

### install

```
go get -u github.com/pbergman/pipe-process
```

### example

```
    cmd, err := pipe.NewProcess(`echo 'SOME SECRET PASSPHRASE' | gpg --batch --decrypt --yes --passphrase-fd 0 dump.sql.gpg | bzip2 -d | mysql example_database`)

    if err != nil {
        panic(err)
    }

    // print command output to default STDOUT
    (*cmd)[cmd.Len()-1].Stdout = os.Stdout

    // map all stderr to default STDERR
    for _, curr := range *cmd {        
        curr.Stderr = os.Stderr		
    }
    
    if err := cmd.Run(); err != nil {
        panic(err) 
    }
```
