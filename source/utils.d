U[V] cdup(U,V)(const U[V] aa)
{
	U[V] bb;
	foreach(k,v; aa)
		bb[k] = v;
	return bb;
}
