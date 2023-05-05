package sqli

var booleanPayloads = []string{
	" OR 1=1",
}

var errPayloads = []string{
	"Syntax error or access violation",
	"Fatal error:",
	"error in your SQL syntax",
	"mysql_num_rows()",
	"mysql_fetch_array()",
	"Error Occurred While Processing Request",
	"Server Error in '/' Application",
	"mysql_fetch_row()",
	"Syntax error",
	"mysql_fetch_assoc()",
	"mysql_fetch_object()",
	"mysql_numrows()",
	"GetArray()",
	"FetchRow()",
	"Input string was not in a correct format",
	"You have an error in your SQL syntax",
	"Warning: session_start()",
	"Warning: is_writable()",
	"Warning: Unknown()",
	"Warning: mysql_result()",
	"Warning: mysql_query()",
	"Warning: mysql_num_rows()",
	"Warning: array_merge()",
	"Warning: preg_match()",
	"SQL syntax error",
	"MYSQL error message: supplied argument….",
	"mysql error with query",
}

var unionPayloads = []string{}

var characters = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "_", "", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "@", ".",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}
