package squel

func argInList(item interface{}, args []interface{}) int {
	for i, arg := range args {
		if item == arg {
			return i
		}
	}
	return -1
}

func countArgs(conds []*Condition) int {
	n := 0
	for _, cond := range conds {
		if len(cond.conditions) > 0 {
			n = countArgs(cond.conditions)
		} else {
			n += len(cond.args)
		}
	}

	return n
}
