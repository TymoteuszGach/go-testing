# go-testing

**Requirements:**

Create an app that reads a JSON message containing two numbers from NATS Streaming, calculates Greatest Common Factor and saves the result to MongoDB. Additionally add a REST endpoint to expose the list of all calculated results.


Problems:
- how to test it? (feedback loop)
- what if requirements change? for example, we would like to store the data in DynamoDB