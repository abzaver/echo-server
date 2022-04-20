/*
Prototype is a creational design pattern that allows cloning objects, even
complex ones, without coupling to their specific classes.

All prototype classes should have a common interface that makes it possible to
copy objects even if their concrete classes are unknown. Prototype objects can
produce full copies since objects of the same class can access each other’s
private fields.

Conceptual Example
Let’s try to figure out the Prototype pattern using an example based on the
operating system’s file system. The OS file system is recursive: the folders
contain files and folders, which may also include files and folders, and so on.

Each file and folder can be represented by an inode interface. inode interface
also has the clone function.

Both file and folder structs implement the print and clone functions since they
are of the inode type. Also, notice the clone function in both file and folder.
The clone function in both of them returns a copy of the respective file or
folder. During the cloning, we append the suffix “_clone” to the name field.
*/

// Client code
package main

import "fmt"

func main() {
	file1 := &File{name: "File1"}
	file2 := &File{name: "File2"}
	file3 := &File{name: "File3"}

	folder1 := &Folder{
		children: []Inode{file1},
		name:     "Folder1",
	}

	folder2 := &Folder{
		children: []Inode{folder1, file2, file3},
		name:     "Folder2",
	}
	fmt.Println("\nPrinting hierarchy for Folder2")
	folder2.print("  ")

	cloneFolder := folder2.clone()
	fmt.Println("\nPrinting hierarchy for clone Folder")
	cloneFolder.print("  ")
}
