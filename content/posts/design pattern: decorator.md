---
title: "Decorator pattern"
date: 2022-11-04
authors: ["ngoctd"]
draft: true
tags: ["design pattern"]
series: ["head first design pattern"]
---

I used to think real men subclassed everything. That was until I learned the power of extension at runtime, rather than at compile time.

Just call this chapter "Design Eye for Inheritance Guy." We'll re-examine the typical overuse of inheritance and you'll learn how to decorate your classes at runtime using a form of object composition. Why? Once you know the techniques of decorating, you'll be able to give your objects new responsibilities without making any code changes to the underlying classes.

## Welcome to Starbuzz Coffee

Starbuzz Coffee has made a name for itself as the fastest-growing coffee shop. When they first went into business they designed their classes like this...

![inheritance version](../../images/design-patterns/decorator/decorator1.png)

In addition to your coffee, you can also ask for several condiments like steamed milk, soy, and mocha...

![inheritance version](../../images/design-patterns/decorator/decorator2.png)

It's pretty obvious that Starbuzz has created a maintenance nightmare for themselves. What happens when the price of milk goes up? What do they do when they add a new caramel topping.

Solution: keep track of the condiments

![inheritance version](../../images/design-patterns/decorator/decorator3.png)

```java
public class Beverage {
    public String description;
    public double milk;
    public double soy;
    public double mocha;

    public double cost() {
        // check condiments here
    }
}

public class DarkRoast extends Beverage {
    public DarkRoast() {
        description = "Dark road description";
    }

    public double cost() {

    }
}
```

A: See, five classes total. This is definitely the way to go.
B: I can see more potential problems with this approach by thinking about how th design might need to change in the future.

- Price changes for condiments will force us to alter existing code.
- New condiments will force us to add new methods and alter the cost method in the superclass
- We may have new beverages. For some of these beverages (tea), the condiments may not be appropriate, yet the Tea subclass will still inherit methods like hasWhip().
- What if a customer wants a double mocha?

## Constructing a drink order with Decorators

- Decorators have the same supertype as the objects they decorate.
- You can use one or more decorators to wrap an object
- Give that the decorator has the same supertype as the object it decorates, we can pass around a decorated object in place of the original wrapped object
- The decorator adds its own behavior before and/or after delegating to the object it decorates to do the rest of the job.
- Objects can be decorated at any time, so we can decorate objects dynamically at runtime as many decorators as we like.

**References**
- Head first design pattern


## Summary

- Thinking beyond the maintenance problem, which of the design principles.

- When i inherit behavior by subclassing, that behavior is set statically at compile time. In addition, all subclasses must inherit the same behavior. However, i can extend an object's behavior through composition, then i can do this dynamically at runtime.

- The Decorator Pattern attaches additional responsibilities to an object dynamically.  Decorators provide a flexible alternatively to subclassing for extending functionality.

129