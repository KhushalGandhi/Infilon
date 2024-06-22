# Infilon

## Description
This project is a sample database schema for managing information about people, their phone numbers, and addresses. It includes tables for storing person details, phone numbers, addresses, and the relationships between people and addresses.

## Database Schema
- **person:** Stores information about individuals, including their name and age.
- **phone:** Stores phone numbers associated with individuals.
- **address:** Stores address details such as city, state, street, and zip code.
- **address_join:** Maps individuals to their addresses.

## Table Structure
1. **person:**
   - id: SERIAL (Primary Key)
   - name: VARCHAR(255)
   - age: INT

2. **phone:**
   - id: SERIAL (Primary Key)
   - number: VARCHAR(255)
   - person_id: INT (Foreign Key referencing person.id)

3. **address:**
   - id: SERIAL (Primary Key)
   - city: VARCHAR(255)
   - state: VARCHAR(255)
   - street1: VARCHAR(255)
   - street2: VARCHAR(255)
   - zip_code: VARCHAR(255)

4. **address_join:**
   - id: SERIAL (Primary Key)
   - person_id: INT (Foreign Key referencing person.id)
   - address_id: INT (Foreign Key referencing address.id)

## Sample Data
- **person:**
  - (1, 'mike', 31)
  - (2, 'John', 20)
  - (3, 'Joseph', 20)

- **phone:**
  - (1, '444-444-4444', 1)
  - (2, '123-444-7777', 2)
  - (3, '445-222-1234', 3)

- **address:**
  - (1, 'Eugene', 'OR', '111 Main St', '', '98765')
  - (2, 'Sacramento', 'CA', '432 First St', 'Apt 1', '22221')
  - (3, 'Austin', 'TX', '213 South 1st St', '', '78704')

- **address_join:**
  - (1, 1, 3)
  - (2, 2, 1)
  - (3, 3, 2)

## Routes
1. **GET /person/{person_id}/info**
   - Returns information about a person including their name, phone number, city, state, street1, street2, and zip code.
   - Example: `GET /person/1/info` returns:
     ```json
     {
         "name": "mike",
         "phone_number": "444-444-4444",
         "city" : "Austin",
         "state" : "TX",
         "street1": "213 South 1st St",
         "street2": "",
         "zip_code": "78704"
     }
     ```

2. **POST /person/create**
   - Creates a new person with the provided details.
   - Request Body:
     ```json
     {
         "name": "YOURNAME",
         "phone_number": "123-456-7890",
         "city" : "Sacramento",
         "state" : "CA",
         "street1": "112 Main St",
         "street2": "Apt 12",
         "zip_code": "12345"
     }
     ```
   - Returns: 200 OK
