# fsaed: Finite State Autonoma EDitor
A tool for editing files according to the rules of a provided Finite State Machine 


## Examples

### Run sed only after seeing multiple patterns

Given the input:

```
baz
foo
baz
bar
baz
```

And you only want to edit the final `baz` into `bang`, use this command:

```
$ echo "baz\nfoo\nbaz\nbar\nbaz" | fsaed '/foo/ /bar/ do s/baz/bang'
```

Results in: 

```
baz
foo
baz
bar
bang
```

### Print Lines Between /regex/

Given the input:

```
DO NOT PRINT THIS LINE
baz - DO NOT PRINT THIS EITHER
foo
bar
baz - DO NOT PRINT THIS EITHER
DO NOT PRINT THIS LINE
```

And you only want to print what's between the `baz`s

```
$ fsaed -n '/baz/ /baz/;print' < file.txt
```

Results In:

```
foo
bar
```

## Syntax

### Statement

### Action
