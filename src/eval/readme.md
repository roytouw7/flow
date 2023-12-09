# Evaluation
The evaluation phase transforms an AST into it's computed representation.

The process transforms an *ast.Program recursively calling Eval on each node.

## Observable
Each identifier is wrapped in an `Observable` *Object* and stored in the environment. This allows for
everytime an identifier is assigned a new value, for NotifyAll to be called with the updated value.
Every other identifier which value (partially) consisting of this identifier then will update its own value
reflecting this change.

```flow
let a = 1;
let b = a;
a = 7;
b; // evaluates to 7 now
```

## Lazy Evaluation of Identifiers
Due to compositions of identifiers with other identifiers or primitives like:
`let b = a + 7;`
we can not save the evaluated value to b. <br />
If we would save the evaluated values we have no way to compute the value of *b* when *a* changes, for example when
a changes from 1 to 2:

### Eager Identifier Evaluation Example
1. `a = 1;`
2. `b = a + 7;`
3. `a + 7` evaluates to `8`
4. `a = 2;`
5. `b` should now become `9`, but the incoming change only informs `a` became `2`
6. we do not know how to update `b` to become `a + 7` because we lost the knowledge of `7` after evaluation 

### Lazy Identifier Evaluation Example
1. `a = 1;`
2. `b = a + 7;`
3. we do not evaluate `a + 7` but store it as the expression
4. `a = 2;`
5. `b;` on usage we now evaluate `b` which holds the expression `a + 7`
6. `a` is now `2` so `a + 7` evaluates to `9`
7. `b` evaluates to `9`

### Self Referencing Stack Overflow
When assigning in a self referencing way like `a = a + 1;` due too lazy evaluation this would result in a stack overflow
on evaluation of value. Because of this self references are eagerly evaluated and substituted for the current value.