// Package translation /*
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package translation

type InfoMessages string

const (
	InfoMessagesUnknown                   InfoMessages = "unknown"
	InfoMessagesDuplicate                 InfoMessages = "duplicate"
	InfoMessagesFileDetected              InfoMessages = "file_detected"
	InfoMessagesComposed                  InfoMessages = "composed"
	InfoMessagesCannotSaveTempFile        InfoMessages = "temp_save_error"
	InfoMessagesCannotWatchDirectory      InfoMessages = "cannot_watch_directory"
	InfoMessagesWatchingDirectory         InfoMessages = "watch_directory"
	InfoMessagesCannotSerialize           InfoMessages = "cannot_serialize"
	InfoMessagesCannotDeSerialize         InfoMessages = "cannot_deserialize"
	InfoMessagesCannotDeleteFile          InfoMessages = "cannot_delete_file"
	InfoMessagesCannotSaveFile            InfoMessages = "cannot_save_file"
	InfoMessagesDeletingFile              InfoMessages = "deleting_file"
	InfoMessagesSavingFile                InfoMessages = "saving_file"
	InfoMessagesProcessingPendingFiles    InfoMessages = "processing_pending_files"
	InfoMessagesCannotProcessPendingFiles InfoMessages = "cannot_process_pending_files"
	InfoMessagesCannotOpenFile            InfoMessages = "cannot_lock_file"
	InfoMessagesCannotLockFile            InfoMessages = "cannot_lock_file"
	InfoMessagesCannotUnLockFile          InfoMessages = "cannot_unlock_file"
	// Add more message constants as needed
)
