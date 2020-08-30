package sql

// // IndexToSQL converts index struct as a sql string.
// // Return true if it's an inline sql.
// func IndexToSQL(config Config, index rel.Index) (string, bool) {
// 	var (
// 		buffer Buffer
// 		typ    = string(index.Type)
// 	)

// 	buffer.WriteString(typ)

// 	if index.Name != "" {
// 		buffer.WriteByte(' ')
// 		buffer.WriteString(Escape(config, index.Name))
// 	}

// 	buffer.WriteString(" (")
// 	for i, col := range index.Columns {
// 		if i > 0 {
// 			buffer.WriteString(", ")
// 		}
// 		buffer.WriteString(Escape(config, col))
// 	}
// 	buffer.WriteString(")")

// 	optionsToSQL(&buffer, index.Options)
// 	return buffer.String(), true
// }

// func optionsToSQL(buffer *Buffer, options string) {
// 	if options == "" {
// 		return
// 	}

// 	buffer.WriteByte(' ')
// 	buffer.WriteString(options)
// }
