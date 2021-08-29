# gofindthem
Gofindthem is a go library that combines a domain specific language (DSL), that is like a logical expresion, with a substring matching engine (implementations of Aho-Corasick at the moment).
Enabling more complex searches with a easy to read expression.
It supports multiple expressions searchs making it an eficient way to "classify" documents according to the expressions that were matched.

This project was conceived in 2019 while I was searching a way to process millions of documents in a more efficient way.
In my searches I found this post [medium flashtext](https://medium.freecodecamp.org/regex-was-taking-5-days-flashtext-does-it-in-15-minutes-55f04411025f) of Vikash Singh.
And discover the Aho-Corasick, which is an awesome algorithm, but it didn't solve my problem completely. 
I needed something that could process as fast as the Aho-Corasick implementations could, but also, that enable more complex searches.
The other problem was that the expressions were feeded by a team of analysts that were used to use regex to classify those documents.
So I needed a syntax that were easy enough to convince them to change.
The idea was to create a DSL that would have the operators "AND", "OR" and "NOT" that were the same of the logical operations and a new operator "INORD"
that would check if the terms were found on the same order that they were specified.
This would allow searches that used the regexe `foo.*bar` to be replaced with `INORD("foo" and "bar")` and the combination of regexes 
`foo.*bar` and `bar.*foo` to be replaced as `"foo" and "bar"`. This is was not a replacement for regex, but it fitted must use cases.
To fill the other cases that only regex would solve, I added a way to represent a regex on the syntax `R"foo.*bar"`.
This way the terms would be searched on the document using an implementations of Aho-Corasick and the regex would use the regex engine.

# Usage