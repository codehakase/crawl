package crawler

import (
	"bufio"
	"bytes"
	"os"
	"sort"
)

func WriteXML(filename, path string, links []string) error {
	sort.Strings(links)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	bufWriter := bufio.NewWriter(file)
	buff := new(bytes.Buffer)
	// write XML data to buffer
	buff.Write(XMLheader)
	buff.Write(XMLroot)
	// write links
	for _, loc := range links {
		buff.Write(XMLpre)
		buff.WriteString(path)
		buff.WriteString(loc)
		buff.Write(XMLcls)
	}
	buff.WriteString("\n</urlset>")
	_, err = bufWriter.Write(buff.Bytes())
	if err != nil {
		return err
	}
	bufWriter.Flush()
	return nil
}
