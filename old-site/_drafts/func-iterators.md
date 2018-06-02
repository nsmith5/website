# Function Iterators

Many programming languages have the ability to co-opt the syntax of a function definition to construct a type of iterator. In Python, for example, instead of the traditional `return` statement, `yield` statements are used to produce the values of the iterator one by one. 

```python
def pascal(n):
	i = 1
	while i <= n:
		yield i
		i += 1
```

While dressed like a function `pascal`, in this case, is actually an iterator called a `generator`. All of the typical iterator mechanics work with it. 

```python
>> p = pascal(10)
>> next(p)
1
>> next(p)
2
>> n = 20
>> sum(pascal(n)) == n * (n + 1) / 2
True
```

