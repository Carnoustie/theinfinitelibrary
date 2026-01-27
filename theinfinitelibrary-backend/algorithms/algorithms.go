package algorithms


// Computes similarity between booktitles to infer equality of books in light of possible incorrectly spelled book titles by users
func LevenshteinDistance(s1 string, s2 string, len_s1 int, len_s2 int) int {
	if len_s1 == 0 {
		return len_s2
	} else if len_s2 == 0 {
		return len_s1
	} else if len_s1 == len_s2 && s1[len_s1-1] == s2[len_s2-1] {
		return LevenshteinDistance(s1, s2, len_s1-1, len_s2-1)
	} else {
		return 1 + min(LevenshteinDistance(s1, s2, len_s1, len_s2-1), min(LevenshteinDistance(s1, s2, len_s1-1, len_s2), LevenshteinDistance(s1, s2, len_s1-1, len_s2-1)))
	}
}
