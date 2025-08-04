# Create a simple python code that is vulnerable to SQL injection

def vulnerable_function():
    API_KEY = "1234567890"
    user_input = input("Enter a username: ")
    query = f"SELECT * FROM users WHERE username = '{user_input}'"
    print(query)
    print(API_KEY)

vulnerable_function()