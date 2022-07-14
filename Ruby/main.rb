require 'fileutils'

require 'filename_increment'

STRUCTURE = {
  images: {
    folder: 'Images',
    extensions: ['png', 'jpg', 'webp', 'svg', 'gif', 'ico', 'jpeg', 'bmp', 'esp', 'jpeg 2000', 'heif', 'bat', 'cgm',
                 'tif', 'tiff', 'eps', 'raw', 'cr2', 'nef', 'orf', 'sr2']
  },
  videos: {
    folder: 'Videos',
    extensions: %w[mp4 mov wmv fly avi mkv flv mpg webm oog m4p m4v qt swf
                   avchd f4v mpeg-2]
  },
  music: {
    folder: 'Music',
    extensions: %w[mp3 aac flac alac wav aiff dsd pcm m4a wma]
  },
  documents: {
    folder: 'Documents',
    extensions: %w[txt doc docx docx odt xls xlsx ppt pptx]
  },
  psd: {
    folder: 'Psd',
    extensions: ['psd']
  },
  pdf: {
    folder: 'Pdf',
    extensions: ['pdf']
  },
  archive: {
    folder: 'Archive',
    extensions: %w[zip rar 7z tar]
  },
  exe: {
    folder: 'Exe',
    extensions: ['exe']
  },
  torrent: {
    folder: 'Torrent',
    extensions: ['torrent']
  }
}

class Sorting
  attr_reader :threads

  def initialize(data, folder_for_sorting = 'Folder for sorting')
    @folder_for_sorting = folder_for_sorting
    @threads = []
    @paths = []
    @data = data
  end

  def by_extension
    create_folder(@folder_for_sorting)
    check_files
    @paths.each do |path|
      in_thread do
        new_path = path[:new_path]
        while File.exist?(new_path)

          new_path = FilenameIncrement.new(new_path).to_s if File.exist?(new_path)
          puts new_path
        end
        FileUtils.mv(path[:old_path], new_path)
        puts "Move #{path[:old_path]} >>> #{path[:new_path]}"
      end
    end
    @threads.map(&:value)
  end

  private

  def in_thread(&block)
    @threads << Thread.new(&block)
  end

  def create_folder(folder)
    pwd = FileUtils.getwd
    folder = File.join(pwd, folder)
    unless Dir.exist?(folder)
      FileUtils.mkdir_p(folder)
      puts "Create folder: #{folder}"
    end
  end

  def get_all_files_in_folder(folder)
    pwd = FileUtils.getwd
    folder = File.join(pwd, folder)
    Dir.glob("#{folder}/**/*")
  end

  def check_files
    all_files = get_all_files_in_folder(@folder_for_sorting)
    all_files.each do |file|
      @data.each do |_key, value|
        value[:extensions].each do |extension|
          file_extname = File.extname(file).delete('.')
          next unless extension.downcase == file_extname.downcase

          file_name = File.basename(file)
          mod_time_year = File.mtime(file).year.to_s
          folder = File.join('Sorted Files', value[:folder], mod_time_year)
          create_folder(folder)
          new_path = File.join(folder, file_name)
          @paths.push({ old_path: file, new_path: })
        end
      end
    end
  end
end

Sorting.new(STRUCTURE).by_extension
