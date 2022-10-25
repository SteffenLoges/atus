package video

func Fragment(sourceFile, outputFile string) error {
	_, err := _exec("mp4fragment", sourceFile, outputFile)
	return err
}
