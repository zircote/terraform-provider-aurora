package aurora

func resourceDataToInt64(d interface{}) int64 {
	i := d.(int)
	return int64(i)
}

func resourceDataToInt32(d interface{}) int32 {
	v := d.(int)
	return int32(v)
}

func resourceDataToInt(d interface{}) int {
	return d.(int)
}

func resourceDataToStringPtr(d interface{}) *string {
	s := d.(string)
	return &s
}

func resourceDataToString(d interface{}) string {
	return d.(string)
}

func resourceDataToBool(d interface{}) bool {
	return d.(bool)
}
