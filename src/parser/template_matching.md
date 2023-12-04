# Template Matching
Template matching is required for parsing expressions or statements which have a starting token which is not unique.

One example is the `(` token. This token can either be a *group expression* or an *arrow function*
```
(1 * (2 + 3)) // grouped expression
const fn = () => { // opening of an arrow function 
```

This can be parse using a template consisting of:
1. regex to match e.g. `\(.+\)\s+=>` for arrow function
2. function to use parsing on match
3. limit, the amount of tokens to peek before returning a negative result for the template

Template Matching starts at the current token, and peeks all tokens up until the limit to find the match?

## Order of Templates
The order of templates is important because matching an *arrow function* will never happen if we 
try matching it agains a *grouped expression* first