#BULK OPERATIONS ON AVL TREES

---
##STRUCTURE
This repository is structured as follows:
    
            avl
                avl.go
                avl_test.go
                utils.go
            cairo-avl
                cairo-avl.go
                cairo_avl_test.go
                dict.go
                node.go
                utils.go
            .gitignore
            main.go
            README.md

The `avl` package contains code relating to the unnested self balancing trees in exact respect to the BFS16 paper.
The `cairo-avl` package contains code relating to the variant implementation contained in the `self balancing tree` document.

The `main.go` file imports both packages for use in comparison of how they work.


##HOW TO RUN
###BUILDING CODE
To run the `main.go` file, in the terminal, from the root directory, run:

    go run main.go

The `main.go` file can be used to run quick tests to observe any cases, feel free to add your code to test specific outcomes.

###RUNNING TESTS
In the `*_test.go` files, there are fuzz tests written for the `union` and `difference` bulk operations. To run any of these fuzz tests,
from the root directory, navigate to the package/folder containing the `*_test.go` file and run the command:

    gotip test -fuzz=FUZZ_FUNCTION_NAME

For example, to run the `FuzzUnion` fuzz test in the `cairo-avl` package, first navigate to the directory from the root directory with:

    > cd cairo-avl

Then run the `FuzzUnion` fuzz tests with:

    > gotip test -fuzz=FuzzUnion

Then you should see the output in the terminal

##COUNTING HASHES
The hashes were counted in two parts, using an unorthodox method.

First, for two trees as inputs to the bulk operation functions, say `t1` and `t2`, all the nodes in `t1` and `t2` were intialised with their
`Exposed` properties set to false. A global variable `numOfExposedNodes` is used to keep a count of all the nodes that have
not been exposed or unrolled yet. The variable is then incremented each time a call to the `exposeNode` function is made.
At the end of the bulk operation, the variable is returned and added to the total hash count of the bulk operation.

The second count kept is the number of nodes created during the lifetime of the bulk operation in question. This count,
as well as the `numOfExposedNodes` variable make up the total hash counts of the operation. This count is gotten by counting
the number of nodes that have their `Exposed` property set to true after the bulk operation (Note that this property is set to true for a new node by default).
