"""
DreamView 2021
Author: Mark Thielen
Date: May 12, 2021

"""
import datetime
import json
import os

# source_path = r'D:\shotgun\source_files'
source_path = r'D:\shotgun\202103'
fix_path = r'D:\shotgun\fix_phage'

def process(sg):
    """process shotgun files

    Args:
        sg (:obj:): shotgun instance
    """
    print('\nSearching for infected files...')
    total_time = None
    total_scan_time = None
    total_file_size = 0

    tags = [{'type': 'Tag', 'id': 4379}]
    any_filters = [
        ['filename', 'ends_with', '.mb'],
        ['filename', 'ends_with', '.ma']]

    filters = [
        # ['id', 'in', [1717073, 1659356, 1073338]],
        # ['created_at', 'greater_than', datetime.datetime(2021, 4, 16)],
        ['created_at', 'in_calendar_month', -2],
        # ['tags', 'not_in', tags],
        {'filter_operator': 'any',
                            'filters': any_filters}
    ]
    fields = list(sg.schema_field_read('Attachment').keys())
    attachments = sg.find('Attachment', filters, fields)

    source_files = {}
    bad = []
    not_bad = []
    time = None

    for attachment in attachments:
        print('Processing: {}'.format(attachment.get('filename')))
        (path, time, file_size) = _download_file(sg, attachment)
        if not time:
            continue

        if not total_time:
            total_time = time
        else:
            total_time = total_time + time
        total_file_size = total_file_size + file_size

        if not path:
            continue

        start_time = datetime.datetime.now()
        suspect = open(path, 'rb').read()
        if b'phage' in suspect:
            end_time = datetime.datetime.now()
            if not total_scan_time:
                total_scan_time = end_time - start_time
            else:
                total_scan_time = total_scan_time + end_time - start_time

            print('****  Bad: {} - Time for search: {}  ****'.format(path, (end_time-start_time)))
            bad.append(path)
            attachment['source_file'] = path
            attachment['source_file_id'] = attachment['id']
            attachment['created_at'] = None
            attachment['updated_at'] = None
            source_files[os.path.basename(path)] = attachment

            with open(os.path.join(source_path, 'source_files.json'), 'a') as f:
                f.write(json.dumps(source_files, indent=4))
            sg.update("Attachment", attachment['id'], {'tags': tags})
        else:
            end_time = datetime.datetime.now()
            print('NOT Bad: {} - Time for search: {}'.format(path, (end_time-start_time)))
            not_bad.append(path)
            try:
                os.remove(path)
            except IOError as e:
                print(e)
                if '[Errno 28] No space left on device' in e:
                    return

    print('\nTotal Size:          {} bytes'.format(total_file_size))
    print('Total Scan Time:     {}'.format(total_scan_time))
    print('Total Download Time: {}\n'.format(total_time))

def _download_file(sg, attachment):
    date = attachment['created_at']

    source_file = 'source_{}{:02d}{:02d}_{}_{}'.format(date.year, date.month, date.day, attachment['id'], attachment.get('filename'))
    source_path_file = os.path.join(source_path, source_file)

    # attachment = {'type': 'Attachment', 'id': attachment_id }
    if os.path.isfile(source_path_file):
        return (None, None, 0)

    start_time = datetime.datetime.now()
    sg.download_attachment(attachment, file_path=source_path_file)
    end_time = datetime.datetime.now()
    file_size = os.path.getsize(source_path_file)
    total_time = end_time - start_time
    print('Download time:  {}\t\t\tFile Size:  {}'.format(total_time, file_size))
    
    return (source_path_file, total_time, file_size)


if __name__ == '__main__':
    process(sg)
