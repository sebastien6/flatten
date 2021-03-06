# flatten
Convert embedded array to a flat array


the package export a single function (FlattenList) that take
an interface as argument and return the resulting flatten
array through an interface.


 FlattenList take as input embedded array in three formats:

 - string: Embedded array inside a string "[1, [2, 3], 4]"
 - []byte: Such as JSON embedded array []byte("[1, [2, 3], 4]")
 - []interface{}: Native Golang embedded array format

 The embedded array is first flatten, and then if all the
 element of the flatten array are of the same type, the
 resulting array is converted in the appropriate array
 format.

 the function can return the flatten array in 4 differents
 format based on its content:

 - []int: 			All elements are integer, not exceeding the limit
 					relative to int32
 - []float64: 		All elements are float value, not exceeding the limit
 					relative to float32
 - []string:    	All elements are of type string or interpreted as string
 - []interface{}:	Element of the array are a mix of different types (integer,
					float and string)

Example:

flatten.FlattenList([]interface{}{5, []interface{}{7, 8}, []interface{}{9, 11, []interface{}{7}}, 1})

Returned: [5 7 8 9 11 7 1], of type []int

flatten.FlattenList([]interface{}{5.1, []interface{}{"hello", 8}})

Returned: [5.1 hello 8], of type []interface {}

flatten.FlattenList([]byte("[1.6, 22.789, [3.1, [4.541], []], 8.96]"))

Returned: [1.6 22.789 3.1 4.541 8.96], of type []float64
