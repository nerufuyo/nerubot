"""
Message constants for the Discord music bot
This file contains all user-facing messages to make them easier to update
"""

# Music player messages
SONG_ADDED_TO_QUEUE = "Added to queue: **{title}**"
NOW_PLAYING_TITLE = "Now playing: **{title}**"
QUEUE_TITLE = "Music Queue"
QUEUE_NOW_PLAYING = "Now Playing:"
QUEUE_UP_NEXT = "Up Next:"
QUEUE_PAGE_INFO = "Page {page}/{max_pages} | {total} songs in queue"
QUEUE_EMPTY = "The queue is empty."
CANNOT_REMOVE_CURRENT = "Cannot remove the currently playing song. Use /skip instead."

# Loop mode messages
LOOP_MODE_INVALID = "Invalid mode. Use 'off', 'song', or 'queue'."
LOOP_MODES = {
    "off": "Loop mode disabled",
    "song": "Looping current song",
    "queue": "Looping entire queue"
}

# Error messages
ERROR_COMMAND = "An error occurred: {error}"
ERROR_NO_RESULTS = "No results found"

# Channel messages
USER_NOT_IN_CHANNEL = "You are not connected to a voice channel."

# Music control messages
SONG_SKIPPED_NEXT = "⏭️ Skipped! Next up: **{title}**"
SONG_SKIPPED_NO_MORE = "⏭️ Skipped! No more songs in queue."
SONG_REMOVED = "Removed **{title}** from the queue."
SONG_REMOVE_INVALID = "Invalid song index."
QUEUE_CLEARED = "Stopped playing and cleared the queue."
QUEUE_SHUFFLED = "Shuffled the queue!"
NOTHING_PLAYING = "I'm not playing anything right now."
NOTHING_PAUSED = "Nothing is paused right now."
PAUSED = "Paused playback."
RESUMED = "Resumed playback."
VOLUME_SET = "Volume set to {volume}%"
VOLUME_RANGE = "Volume must be between 0 and 100."

# Music service messages
PLAYLIST_NO_VIDEOS = "❌ No playable videos found in this playlist."
PLAYLIST_FOUND = "🎵 Found playlist with {count} available videos. Playing first video: **{title}**"
NO_AUDIO_STREAM = "❌ Could not find a playable audio stream for this video."
NOW_PLAYING_NOTIFICATION = "Now playing: **{title}**"

# Help command descriptions
HELP_JOIN_DESC = "Join your voice channel"
HELP_LEAVE_DESC = "Leave the voice channel"
HELP_PLAY_DESC = "Play a song from YouTube, Spotify, or SoundCloud"
HELP_STOP_DESC = "Stop music and clear queue"
HELP_PAUSE_DESC = "Pause the current song"
HELP_RESUME_DESC = "Resume the current song"
HELP_SKIP_DESC = "Skip the current song"
HELP_VOLUME_DESC = "Set the volume (0-100)"
HELP_NOW_DESC = "Show current song with status"
HELP_QUEUE_DESC = "Show the music queue"
HELP_REMOVE_DESC = "Remove a song from the queue"
HELP_SHUFFLE_DESC = "Shuffle the music queue"
HELP_LOOP_DESC = "Toggle loop mode (off/single/queue)"