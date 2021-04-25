# Distributed Word Counter
- The input file is divided evenly among 5 go routines, each routine computes the word counts for the portion of the file it is responsible for. 
- After each routine finishes, it writes the output to a shared map that is handled by another routine(the reducer). 
- Only one routine can access the map at a time and the output correctness is guaranteed.
- When all routines write the output to the shared map, the “reducer” writes the output in the file sorted by the frequency.
